package genose

import (
	"errors"
	"fmt"
	"github.com/genose-projects/goo"
	configure "github.com/genose-projects/genose-configure"
	context "github.com/genose-projects/genose-context"
	core "github.com/genose-projects/genose-core"
	web "github.com/genose-projects/genose-web"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

var bannerText = []string{"",
"         ",
"                  888b    888",                           
"                  8888b   888",                            
"                  88888b  888",                            
".d88b.   .d88b.   888Y88b 888  .d88b.  .d8888b   .d88b.",
"d88P`88b d8P  Y8b 888 Y88b888 d88``88b 88K      d8P  Y8b", 
"888  888 88888888 888  Y88888 888  888 `Y8888b. 88888888 ",
"Y88b 888 Y8b.     888   Y8888 Y88..88P      X88 Y8b.",     
"  Y88888 `Y8888   888    Y888  `Y88P`   88888P'  `Y8888`" , 
"     888",                                                 
"Y8b d88P",                                                 
"  Y88P",  
"      ",
}

type application interface {
	getLogger() context.Logger
	getLoggingProperties(arguments ApplicationArguments) *configure.LoggingProperties
	configureLogger(logger context.Logger, loggingProperties *configure.LoggingProperties)
	getTaskWatch() *core.TaskWatch
	getApplicationId() context.ApplicationId
	getContextId() context.ContextId
	generateApplicationAndContextId()
	getApplicationArguments() ApplicationArguments
	printBanner()
	logStarting()
	scanComponents(arguments ApplicationArguments) error
	prepareEnvironment(arguments ApplicationArguments, listeners *ApplicationRunListeners) (core.Environment, error)
	prepareContext(environment core.ConfigurableEnvironment, arguments ApplicationArguments, listeners *ApplicationRunListeners, loggingProperties *configure.LoggingProperties) (context.ConfigurableApplicationContext, error)
	getApplicationRunListenerInstances(arguments ApplicationArguments) (*ApplicationRunListeners, error)
	getApplicationListeners() []context.ApplicationListener
	getApplicationContextInitializers() []context.ApplicationContextInitializer
	initApplicationListenerInstances() error
	initApplicationContextInitializers() error
	logStarted()
	invokeApplicationRunners(ctx context.ApplicationContext, arguments ApplicationArguments)
	finish()
}

type environmentProvider interface {
	getNewEnvironment() core.ConfigurableEnvironment
}

type contextProvider interface {
	getNewContext(applicationId context.ApplicationId, contextId context.ContextId) context.ConfigurableApplicationContext
}

type GenoseApplication struct {
	application
}

func NewGenoseApplication() *GenoseApplication {
	baseApplication := newBaseApplication()
	app := &GenoseApplication{
		baseApplication,
	}
	baseApplication.genoseApplication = app
	return app
}

func (genoseApp *genoseApplication) Run() {
	taskWatch := genoseApp.getTaskWatch()
	taskWatch.Start()

	// get the application arguments
	arguments := genoseApp.getApplicationArguments()

	logger := genoseApp.getLogger()
	loggingProperties := genoseApp.getLoggingProperties(arguments)
	genoseApp.configureLogger(logger, loggingProperties)

	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				err := r.(error)
				errorString := err.Error()
				logger.Fatal(genoseApp.getContextId(), errorString+"\n"+string(debug.Stack()))
			case string:
				errorString := r.(string)
				logger.Fatal(genoseApp.getContextId(), errorString+"\n"+string(debug.Stack()))
			default:
				logger.Error(genoseApp.getContextId(), r)
				logger.Fatal(genoseApp.getContextId(), string(debug.Stack()))
			}
		}
	}()

	genoseApp.printBanner()

	// log starting
	genoseApp.logStarting()

	// scan components
	err := genoseApp.scanComponents(arguments)
	if err != nil {
		logger.Fatal(genoseApp.getContextId(), err)
	}

	// application listener
	err = genoseApp.initApplicationListenerInstances()
	if err != nil {
		logger.Fatal(genoseApp.getContextId(), err)
	}

	// application context initializers
	err = genoseApp.initApplicationContextInitializers()
	if err != nil {
		logger.Fatal(genoseApp.getContextId(), err)
	}

	// app run listeners
	var listeners *ApplicationRunListeners
	listeners, err = genoseApp.getApplicationRunListenerInstances(arguments)
	if err != nil {
		logger.Fatal(genoseApp.getContextId(), err)
	}

	// broadcast an event to inform the application is starting
	listeners.OnApplicationStarting()

	// prepare environment
	var environment core.Environment
	environment, err = genoseApp.prepareEnvironment(arguments, listeners)
	if err != nil {
		logger.Fatal(genoseApp.getContextId(), err)
	}

	// prepare context
	var applicationContext context.ConfigurableApplicationContext
	applicationContext, err = genoseApp.prepareContext(environment.(core.ConfigurableEnvironment),
		arguments,
		listeners,
		loggingProperties,
	)
	if err != nil {
		logger.Fatal(genoseApp.getContextId(), err)
	}
	taskWatch.Stop()
	genoseApp.logStarted()

	listeners.OnApplicationStarted(applicationContext)
	genoseApp.invokeApplicationRunners(applicationContext, arguments)
	listeners.OnApplicationRunning(applicationContext)

	genoseApp.finish()
}

type baseApplication struct {
	genoseApplication   *GenoseApplication
	applicationId       context.ApplicationId
	contextId           context.ContextId
	logger              context.Logger
	customLogger        context.Logger
	taskWatch           *core.TaskWatch
	listeners           []context.ApplicationListener
	contextInitializers []context.ApplicationContextInitializer
	contextProvider     contextProvider
	environmentProvider environmentProvider
}

func newBaseApplication() *baseApplication {
	baseApplication := &baseApplication{
		listeners:           make([]context.ApplicationListener, 0),
		contextInitializers: make([]context.ApplicationContextInitializer, 0),
		taskWatch:           core.NewTaskWatch(),
		logger:              context.NewSimpleLogger(),
		contextProvider:     newDefaultContextProvider(),
		environmentProvider: newDefaultEnvironmentProvider(),
	}
	baseApplication.generateApplicationAndContextId()

	return baseApplication
}

func (application *baseApplication) getLogger() context.Logger {
	if application.customLogger != nil {
		return application.customLogger
	}
	return application.logger
}

func (application *baseApplication) getTaskWatch() *core.TaskWatch {
	return application.taskWatch
}

func (application *baseApplication) getApplicationId() context.ApplicationId {
	return application.applicationId
}

func (application *baseApplication) getContextId() context.ContextId {
	return application.contextId
}

func (application *baseApplication) printBanner() {
	logger := application.getLogger()
	for _, line := range bannerText {
		logger.Print(application.contextId, line)
	}
}

func (application *baseApplication) getApplicationArguments() ApplicationArguments {
	return getApplicationArguments(os.Args)
}

func (application *baseApplication) generateApplicationAndContextId() {
	var applicationId [36]byte
	core.GenerateUUID(applicationId[:])
	var contextId [36]byte
	core.GenerateUUID(contextId[:])

	application.applicationId = context.ApplicationId(applicationId[:])
	application.contextId = context.ContextId(contextId[:])
}

func (application *baseApplication) prepareEnvironment(arguments ApplicationArguments, listeners *ApplicationRunListeners) (core.Environment, error) {
	environment := application.environmentProvider.getNewEnvironment()

	propertySources := environment.GetPropertySources()
	if arguments != nil && len(arguments.GetSourceArgs()) > 0 {
		propertySources.Add(core.NewSimpleCommandLinePropertySource(arguments.GetSourceArgs()))
	}

	propertySources.Add(core.NewSystemEnvironmentPropertySource())

	listeners.OnApplicationEnvironmentPrepared(environment)
	return environment, nil
}

func (application *baseApplication) scanComponents(arguments ApplicationArguments) error {
	if arguments == nil {
		return nil
	}
	argumentComponentScan := arguments.GetOptionValues("genose.component.scan")
	if argumentComponentScan != nil && len(argumentComponentScan) == 1 && argumentComponentScan[0] == "false" {
		return nil
	}

	application.logger.Info(application.contextId, "Scanning components...")
	componentScanner := newComponentScanner()
	componentCount, err := componentScanner.scan(application.contextId, application.logger)
	if err != nil {
		return err
	}

	application.logger.Info(application.contextId, fmt.Sprintf("Found (%d) components.", componentCount))
	return nil
}

func (application *baseApplication) prepareContext(environment core.ConfigurableEnvironment,
	arguments ApplicationArguments,
	listeners *ApplicationRunListeners,
	loggingProperties *configure.LoggingProperties) (context.ConfigurableApplicationContext, error) {

	applicationContext := application.contextProvider.getNewContext(application.applicationId, application.contextId)

	if applicationContext == nil {
		return nil, errors.New("context could not be created")
	}

	// set environment
	applicationContext.SetEnvironment(environment)
	// set logger
	applicationContext.SetLogger(application.logger)
	factory := applicationContext.GetPeaFactory()

	// apply context initializers
	for _, contextInitializer := range application.getApplicationContextInitializers() {
		contextInitializer.InitializeContext(applicationContext)
	}
	factory.ExcludeType(goo.GetType((*ApplicationRunListener)(nil)))

	// broadcast an event to notify that context is prepared
	listeners.OnApplicationContextPrepared(applicationContext)

	// register application arguments as shared pea
	err := factory.RegisterSharedPea("genoseApplicationArguments", arguments)
	if err != nil {
		return nil, err
	}

	err = factory.RegisterSharedPea("loggingProperties", loggingProperties)
	if err != nil {
		return nil, err
	}

	// broadcast an event to notify that context is loaded
	listeners.OnApplicationContextLoaded(applicationContext)

	if configurableContextAdapter, ok := applicationContext.(context.ConfigurableContextAdapter); ok {
		configurableContextAdapter.Configure()
		return applicationContext, nil
	}
	return nil, errors.New("context.ConfigurableContextAdapter methods must be implemented in your context struct")
}

func (application *baseApplication) getApplicationRunListenerInstances(arguments ApplicationArguments) (*ApplicationRunListeners, error) {
	instances, err := getInstancesWithParamTypes(goo.GetType((*ApplicationRunListener)(nil)),
		[]goo.Type{goo.GetType((*GenoseApplication)(nil)), goo.GetType((*ApplicationArguments)(nil))},
		[]interface{}{application.genoseApplication, arguments})
	if err != nil {
		return nil, err
	}
	var listeners []ApplicationRunListener
	for _, instance := range instances {
		listeners = append(listeners, instance.(ApplicationRunListener))
	}
	return NewApplicationRunListeners(listeners), nil
}

func (application *baseApplication) getApplicationListeners() []context.ApplicationListener {
	return application.listeners
}

func (application *baseApplication) getApplicationContextInitializers() []context.ApplicationContextInitializer {
	return application.contextInitializers
}

func (application *baseApplication) initApplicationListenerInstances() error {
	instances, err := getInstances(goo.GetType((*context.ApplicationListener)(nil)))
	if err != nil {
		return err
	}
	listenerInstances := make([]context.ApplicationListener, len(instances))
	for index, instance := range instances {
		listenerInstances[index] = instance.(context.ApplicationListener)
	}
	application.listeners = listenerInstances
	return nil
}

func (application *baseApplication) initApplicationContextInitializers() error {
	instances, err := getInstances(goo.GetType((*context.ApplicationContextInitializer)(nil)))
	if err != nil {
		return err
	}
	initializerInstances := make([]context.ApplicationContextInitializer, len(instances))
	for index, instance := range instances {
		initializerInstances[index] = instance.(context.ApplicationContextInitializer)
	}
	application.contextInitializers = initializerInstances
	return nil
}

func (application *baseApplication) invokeApplicationRunners(ctx context.ApplicationContext, arguments ApplicationArguments) {
	applicationRunners := ctx.GetSharedPeasByType(goo.GetType((*ApplicationRunner)(nil)))
	for _, applicationRunner := range applicationRunners {
		applicationRunner.(ApplicationRunner).OnApplicationRun(ctx, arguments)
	}
}

func (application *baseApplication) logStarting() {
	logger := application.logger
	if application.customLogger != nil {
		logger = application.customLogger
	}

	logger.Info(application.contextId, "Starting...")
	logger.Infof(application.contextId, "Application Id : %s", application.applicationId)
	logger.Infof(application.contextId, "Application Context Id : %s", application.contextId)
	logger.Info(application.contextId, "Running with Genose, Genose "+Version)
}

func (application *baseApplication) logStarted() {
	lastTime := float32(application.taskWatch.GetTotalTime()) / 1e9
	formattedText := fmt.Sprintf("Started in %.2f second(s)", lastTime)
	application.logger.Info(application.contextId, formattedText)
}

func (application *baseApplication) finish() {
	exitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(exitSignalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignalChannel
}

func (application *baseApplication) getCustomLogger() {
	customLoggers, err := getInstances(goo.GetType((*context.Logger)(nil)))
	if err != nil {
		panic(err)
	}

	if customLoggers != nil {
		if len(customLoggers) != 1 {
			panic("Custom logger cannot be distinguished because there are more than one")
		}

		if len(customLoggers) != 0 {
			application.customLogger = customLoggers[0].(context.Logger)
		}
	}
}

func (application *baseApplication) configureLogger(logger context.Logger, loggingProperties *configure.LoggingProperties) {
	if logger == nil {
		return
	}

	if configurableLogger, ok := logger.(context.LoggingConfiguration); ok {
		configurableLogger.ApplyLoggingProperties(*loggingProperties)
	}
}

func (application *baseApplication) getLoggingProperties(arguments ApplicationArguments) *configure.LoggingProperties {
	if arguments == nil {
		return nil
	}

	properties := &configure.LoggingProperties{}
	loggingLevel := arguments.GetOptionValues("logging.level")
	if len(loggingLevel) != 0 {
		properties.Level = loggingLevel[0]
	} else {
		properties.Level = "DEBUG"
	}

	loggingFile := arguments.GetOptionValues("logging.file.name")
	if len(loggingFile) != 0 {
		properties.FileName = loggingFile[0]
	}

	loggingPath := arguments.GetOptionValues("logging.file.path")
	if len(loggingPath) != 0 {
		properties.FilePath = loggingPath[0]
	}

	return properties
}

type defaultEnvironmentProvider struct {
}

func newDefaultEnvironmentProvider() defaultEnvironmentProvider {
	return defaultEnvironmentProvider{}
}

func (provider defaultEnvironmentProvider) getNewEnvironment() core.ConfigurableEnvironment {
	return web.NewStandardWebEnvironment()
}

type defaultContextProvider struct {
}

func newDefaultContextProvider() defaultContextProvider {
	return defaultContextProvider{}
}

func (provider defaultContextProvider) getNewContext(applicationId context.ApplicationId, contextId context.ContextId) context.ConfigurableApplicationContext {
	return web.NewGenoseServerApplicationContext(applicationId, contextId)
}
