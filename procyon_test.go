package genose

import (
	"errors"
	"fmt"
	"github.com/genose-projects/goo"
	configure "github.com/genose-projects/genose-configure"
	context "github.com/genose-projects/genose-context"
	core "github.com/genose-projects/genose-core"
	peas "github.com/genose-projects/genose-peas"
	web "github.com/genose-projects/genose-web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
	"time"
)

func init() {
	core.Register(newTestApplicationContextInitializer)
}

type applicationMock struct {
	mock.Mock
}

func (app *applicationMock) getLogger() context.Logger {
	results := app.Called()
	return results.Get(0).(context.Logger)
}

func (app *applicationMock) getLoggingProperties(arguments ApplicationArguments) *configure.LoggingProperties {
	results := app.Called(arguments)
	return results.Get(0).(*configure.LoggingProperties)
}

func (app *applicationMock) configureLogger(logger context.Logger, loggingProperties *configure.LoggingProperties) {
	//app.Called(logger, loggingProperties)
}

func (app *applicationMock) getTaskWatch() *core.TaskWatch {
	results := app.Called()
	return results.Get(0).(*core.TaskWatch)
}

func (app *applicationMock) getApplicationId() context.ApplicationId {
	results := app.Called()
	return results.Get(0).(context.ApplicationId)
}

func (app *applicationMock) getContextId() context.ContextId {
	results := app.Called()
	return results.Get(0).(context.ContextId)
}

func (app *applicationMock) printBanner() {
	app.Called()
}

func (app *applicationMock) getApplicationArguments() ApplicationArguments {
	results := app.Called()
	return results.Get(0).(ApplicationArguments)
}

func (app *applicationMock) generateApplicationAndContextId() {
	app.Called()
}

func (app *applicationMock) prepareEnvironment(arguments ApplicationArguments, listeners *ApplicationRunListeners) (core.Environment, error) {
	results := app.Called(arguments, listeners)
	return results.Get(0).(core.Environment), results.Error(1)
}

func (app *applicationMock) scanComponents(arguments ApplicationArguments) error {
	results := app.Called(arguments)
	return results.Error(0)
}

func (app *applicationMock) prepareContext(environment core.ConfigurableEnvironment,
	arguments ApplicationArguments,
	listeners *ApplicationRunListeners,
	loggingProperties *configure.LoggingProperties) (context.ConfigurableApplicationContext, error) {
	results := app.Called(environment, arguments, listeners, loggingProperties)
	return results.Get(0).(context.ConfigurableApplicationContext), results.Error(1)
}

func (app *applicationMock) getApplicationRunListenerInstances(arguments ApplicationArguments) (*ApplicationRunListeners, error) {
	results := app.Called(arguments)
	return results.Get(0).(*ApplicationRunListeners), results.Error(1)
}

func (app *applicationMock) getApplicationListeners() []context.ApplicationListener {
	results := app.Called()
	return results.Get(0).([]context.ApplicationListener)
}

func (app *applicationMock) getApplicationContextInitializers() []context.ApplicationContextInitializer {
	results := app.Called()
	return results.Get(0).([]context.ApplicationContextInitializer)
}

func (app *applicationMock) initApplicationListenerInstances() error {
	results := app.Called()
	return results.Error(0)
}

func (app *applicationMock) initApplicationContextInitializers() error {
	results := app.Called()
	return results.Error(0)
}

func (app *applicationMock) invokeApplicationRunners(ctx context.ApplicationContext, arguments ApplicationArguments) {
	app.Called(ctx, arguments)
}

func (app *applicationMock) logStarting() {
	app.Called()
}

func (app *applicationMock) logStarted() {
	app.Called()
}

func (app *applicationMock) finish() {
	app.Called()
}

func TestGenoseApplication_NewGenoseApplication(t *testing.T) {
	genoseApp := NewGenoseApplication()
	assert.NotNil(t, genoseApp.getContextId())
	assert.NotNil(t, genoseApp.getApplicationId())
}

func TestGenoseApplication_Run_Successfully(t *testing.T) {
	loggingProperties := &configure.LoggingProperties{}
	var applicationIdArray [36]byte
	core.GenerateUUID(applicationIdArray[:])
	var contextIdArray [36]byte
	core.GenerateUUID(contextIdArray[:])

	logger := context.NewSimpleLogger()
	taskWatch := core.NewTaskWatch()
	applicationRunListeners := NewApplicationRunListeners(nil)

	genoseApplication := NewGenoseApplication()

	contextId := context.ContextId(contextIdArray[:])
	applicationId := context.ApplicationId(applicationIdArray[:])

	mockApplication := &applicationMock{}
	genoseApplication.application = mockApplication

	mockApplication.On("getLoggingProperties", mock.Anything).Return(loggingProperties)
	mockApplication.On("getLogger").Return(logger)
	mockApplication.On("getTaskWatch").Return(taskWatch)
	//mockApplication.On("getApplicationId").Return(baseApplication.applicationId)
	//mockApplication.On("getContextId").Return(baseApplication.contextId)

	mockApplication.On("printBanner")
	mockApplication.On("logStarting")

	applicationArguments := getApplicationArguments(nil)
	mockApplication.On("getApplicationArguments").Return(applicationArguments)

	mockApplication.On("scanComponents", applicationArguments).Return(nil)

	mockApplication.On("initApplicationListenerInstances").Return(nil)

	mockApplication.On("initApplicationContextInitializers").Return(nil)

	mockApplication.On("getApplicationRunListenerInstances", applicationArguments).
		Return(applicationRunListeners, nil)

	environment := web.NewStandardWebEnvironment()
	mockApplication.On("prepareEnvironment", applicationArguments, applicationRunListeners).
		Return(environment, nil)

	applicationContext := web.NewGenoseServerApplicationContext(applicationId, contextId)
	mockApplication.On("prepareContext", environment, applicationArguments, applicationRunListeners, loggingProperties).
		Return(applicationContext, nil)

	mockApplication.On("logStarted")

	mockApplication.On("invokeApplicationRunners", applicationContext, applicationArguments)

	mockApplication.On("finish")

	genoseApplication.Run()
	mockApplication.AssertExpectations(t)
}

func TestGenoseApplication_Run_Failed(t *testing.T) {
	loggingProperties := &configure.LoggingProperties{}

	var applicationIdArray [36]byte
	core.GenerateUUID(applicationIdArray[:])
	var contextIdArray [36]byte
	core.GenerateUUID(contextIdArray[:])

	contextId := context.ContextId(contextIdArray[:])
	applicationId := context.ApplicationId(applicationIdArray[:])

	err := errors.New("test error")

	loggerMock := loggerMock{}
	loggerMock.On("Fatal", contextId, err)

	taskWatch := core.NewTaskWatch()
	applicationRunListeners := NewApplicationRunListeners(nil)

	genoseApplication := NewGenoseApplication()

	mockApplication := &applicationMock{}
	mockApplication.On("getContextId").Return(contextId)

	genoseApplication.application = mockApplication

	mockApplication.On("getLoggingProperties", mock.Anything).Return(loggingProperties)
	mockApplication.On("getLogger").Return(loggerMock)
	mockApplication.On("getTaskWatch").Return(taskWatch)

	mockApplication.On("printBanner")
	mockApplication.On("logStarting")

	applicationArguments := getApplicationArguments(nil)
	mockApplication.On("getApplicationArguments").Return(applicationArguments)

	mockApplication.On("scanComponents", applicationArguments).Return(err)

	mockApplication.On("initApplicationListenerInstances").Return(err)

	mockApplication.On("initApplicationContextInitializers").Return(err)

	mockApplication.On("getApplicationRunListenerInstances", applicationArguments).
		Return(applicationRunListeners, err)

	environment := web.NewStandardWebEnvironment()
	mockApplication.On("prepareEnvironment", applicationArguments, applicationRunListeners).
		Return(environment, err)

	applicationContext := web.NewGenoseServerApplicationContext(applicationId, contextId)
	mockApplication.On("prepareContext", environment, applicationArguments, applicationRunListeners, loggingProperties).
		Return(applicationContext, err)

	mockApplication.On("logStarted")

	mockApplication.On("invokeApplicationRunners", applicationContext, applicationArguments)

	mockApplication.On("finish")

	genoseApplication.Run()

	loggerMock.AssertExpectations(t)
	mockApplication.AssertExpectations(t)
}

func TestBaseApplication_getLogger(t *testing.T) {
	assert.NotNil(t, newBaseApplication().getLogger())
}

func TestBaseApplication_getTaskWatch(t *testing.T) {
	assert.NotNil(t, newBaseApplication().getTaskWatch())
}

func TestBaseApplication_getApplicationId(t *testing.T) {
	assert.NotNil(t, newBaseApplication().getApplicationId())
}

func TestBaseApplication_getContextId(t *testing.T) {
	assert.NotNil(t, newBaseApplication().getContextId())
}

func TestBaseApplication_generateApplicationAndContextId(t *testing.T) {
	baseApp := newBaseApplication()
	assert.NotNil(t, baseApp.getContextId())
	assert.NotNil(t, baseApp.getApplicationId())
}

func TestBaseApplication_getApplicationArguments(t *testing.T) {
	assert.NotNil(t, newBaseApplication().getApplicationArguments())
}

func TestBaseApplication_printBanner(t *testing.T) {
	newBaseApplication().printBanner()
}

type loggerMock struct {
	mock.Mock
}

func (l loggerMock) Trace(ctx interface{}, message interface{}) {
}

func (l loggerMock) Debug(ctx interface{}, message interface{}) {
}

func (l loggerMock) Info(ctx interface{}, message interface{}) {
	l.Called(ctx, message)
}

func (l loggerMock) Warning(ctx interface{}, message interface{}) {
}

func (l loggerMock) Error(ctx interface{}, message interface{}) {
}

func (l loggerMock) Fatal(ctx interface{}, message interface{}) {
	l.Called(ctx, message)
}

func (l loggerMock) Panic(ctx interface{}, message interface{}) {
}

func (l loggerMock) Print(ctx interface{}, message interface{}) {
}

func (l loggerMock) Tracef(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Debugf(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Infof(ctx interface{}, format string, args ...interface{}) {
	l.Called(ctx, format, args)
}

func (l loggerMock) Warningf(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Errorf(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Fatalf(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Panicf(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Printf(ctx interface{}, format string, args ...interface{}) {
}

func TestTestBaseApplication_logStarting(t *testing.T) {
	loggerMock := loggerMock{}
	baseApplication := newBaseApplication()
	loggerMock.On("Info", baseApplication.contextId, "Starting...")
	var args = make([]interface{}, 0)
	args = append(args, baseApplication.applicationId)
	loggerMock.On("Infof", baseApplication.contextId, "Application Id : %s", args)

	args = make([]interface{}, 0)
	args = append(args, baseApplication.contextId)
	loggerMock.On("Infof", baseApplication.contextId, "Application Context Id : %s", args)

	loggerMock.On("Info", baseApplication.contextId, "Running with Genose, Genose "+Version)

	baseApplication.logger = loggerMock

	baseApplication.logStarting()
	loggerMock.AssertExpectations(t)
}

func TestBaseApplication_scanComponents(t *testing.T) {
	baseApplication := newBaseApplication()
	baseApplication.scanComponents(getApplicationArguments(os.Args))
}

type environmentProviderMock struct {
	environment environmentMock
}

func newEnvironmentProviderMock(environment environmentMock) environmentProviderMock {
	return environmentProviderMock{
		environment,
	}
}

func (provider environmentProviderMock) getNewEnvironment() core.ConfigurableEnvironment {
	return provider.environment
}

type environmentMock struct {
	core.Environment
	mock.Mock
}

func newEnvironmentMock() environmentMock {
	return environmentMock{}
}

func (env environmentMock) GetPropertySources() *core.PropertySources {
	result := env.Called()
	return result.Get(0).(*core.PropertySources)
}

func (env environmentMock) GetSystemEnvironment() []string {
	result := env.Called()
	return result.Get(0).([]string)
}

func (env environmentMock) GetTypeConverterService() core.TypeConverterService {
	result := env.Called()
	return result.Get(0).(core.TypeConverterService)
}

func TestBaseApplication_prepareEnvironment(t *testing.T) {
	propertySources := core.NewPropertySources()
	environmentMock := newEnvironmentMock()
	environmentMock.On("GetPropertySources").Return(propertySources)

	environmentProviderMock := newEnvironmentProviderMock(environmentMock)

	baseApplication := newBaseApplication()
	baseApplication.environmentProvider = environmentProviderMock
	applicationRunListeners := NewApplicationRunListeners(nil)

	env, err := baseApplication.prepareEnvironment(getApplicationArguments(os.Args), applicationRunListeners)
	assert.NotNil(t, env)
	assert.Nil(t, err)

	assert.Equal(t, 2, propertySources.GetSize())
	environmentMock.AssertExpectations(t)
}

type contextProviderMock struct {
	context contextMock
}

func newContextProviderMock(context contextMock) contextProviderMock {
	return contextProviderMock{
		context,
	}
}

func (provider contextProviderMock) getNewContext(applicationId context.ApplicationId, contextId context.ContextId) context.ConfigurableApplicationContext {
	return provider.context
}

type contextMock struct {
	context.ApplicationContext
	mock.Mock
}

func newContextMock() contextMock {
	return contextMock{}
}

func (ctx contextMock) SetLogger(logger context.Logger) {
	ctx.Called(logger)
}

func (ctx contextMock) GetLogger() context.Logger {
	results := ctx.Called()
	return results.Get(0).(context.Logger)
}

func (ctx contextMock) SetEnvironment(environment core.ConfigurableEnvironment) {
	ctx.Called(environment)
}

func (ctx contextMock) GetEnvironment() core.ConfigurableEnvironment {
	results := ctx.Called()
	return results.Get(0).(core.ConfigurableEnvironment)
}

func (ctx contextMock) GetPeaFactory() peas.ConfigurablePeaFactory {
	results := ctx.Called()
	return results.Get(0).(peas.ConfigurablePeaFactory)
}

func (ctx contextMock) AddApplicationListener(listener context.ApplicationListener) {
	ctx.Called(listener)
}

func (ctx contextMock) Configure() {
	ctx.Called()
}

func (ctx contextMock) OnConfigure() {
	ctx.Called()
}

func (ctx contextMock) FinishConfigure() {
	ctx.Called()
}

type peaFactoryMock struct {
	mock.Mock
	peas.ConfigurablePeaFactory
}

func newPeaFactoryMock() peaFactoryMock {
	return peaFactoryMock{}
}

func (peaFactory peaFactoryMock) RegisterSharedPea(peaName string, sharedObject interface{}) error {
	results := peaFactory.Called(peaName, sharedObject)
	return results.Error(0)
}

func (peaFactory peaFactoryMock) ExcludeType(typ goo.Type) error {
	results := peaFactory.Called(typ)
	return results.Error(0)
}

func TestBaseApplication_prepareContext(t *testing.T) {
	loggingProperties := &configure.LoggingProperties{}

	arguments := getApplicationArguments(os.Args)
	peaFactoryMock := newPeaFactoryMock()
	peaFactoryMock.On("RegisterSharedPea", "genoseApplicationArguments", arguments).Return(nil)
	peaFactoryMock.On("RegisterSharedPea", "loggingProperties", loggingProperties).Return(nil)
	peaFactoryMock.On("ExcludeType", mock.Anything).Return(nil)

	environment := web.NewStandardWebEnvironment()
	contextMock := newContextMock()

	baseApplication := newBaseApplication()

	contextMock.On("SetEnvironment", environment)
	contextMock.On("SetLogger", baseApplication.logger)
	contextMock.On("Configure")
	contextMock.On("GetPeaFactory").Return(peaFactoryMock)
	contextProviderMock := newContextProviderMock(contextMock)

	baseApplication.contextProvider = contextProviderMock
	applicationRunListeners := NewApplicationRunListeners(nil)
	ctx, err := baseApplication.prepareContext(environment, arguments, applicationRunListeners, loggingProperties)

	assert.NotNil(t, ctx)
	assert.Nil(t, err)

	contextMock.AssertExpectations(t)
	peaFactoryMock.AssertExpectations(t)
}

func TestBaseApplication_getAppRunListenerInstances(t *testing.T) {
	genoseApp := NewGenoseApplication()
	runListeners, err := genoseApp.getApplicationRunListenerInstances(getApplicationArguments(os.Args))
	assert.NotNil(t, runListeners)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(runListeners.listeners))
}

type testApplicationContextInitializer struct {
}

func newTestApplicationContextInitializer() testApplicationContextInitializer {
	return testApplicationContextInitializer{}
}

func (initializer testApplicationContextInitializer) InitializeContext(context context.ConfigurableApplicationContext) {

}

func TestBaseApplication_initApplicationContextInitializers(t *testing.T) {
	genoseApp := NewGenoseApplication()
	err := genoseApp.initApplicationContextInitializers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(genoseApp.getApplicationContextInitializers()))
}

func TestTestBaseApplication_logStarted(t *testing.T) {
	loggerMock := loggerMock{}
	baseApplication := newBaseApplication()

	taskWatch := core.NewTaskWatch()
	taskWatch.Start()
	time.Sleep(1000)
	taskWatch.Stop()

	lastTime := float32(taskWatch.GetTotalTime()) / 1e9
	formattedText := fmt.Sprintf("Started in %.2f second(s)", lastTime)
	loggerMock.On("Info", baseApplication.contextId, formattedText)

	baseApplication.logger = loggerMock
	baseApplication.logStarted()
	loggerMock.AssertExpectations(t)
}
