package genose

import (
	context "github.com/genose-projects/genose-context"
	core "github.com/genose-projects/genose-core"
	"time"
)

var applicationEventId = context.GetEventId("github.com.genose.ApplicationEventId")
var applicationStartingEventId = context.GetEventId("github.com.genose.ApplicationStartingEvent")
var applicationEnvironmentPreparedEventId = context.GetEventId("github.com.genose.ApplicationEnvironmentPreparedEvent")
var applicationContextInitializedEventId = context.GetEventId("github.com.genose.ApplicationContextInitializedEvent")
var applicationPreparedEventId = context.GetEventId("github.com.genose.ApplicationPreparedEvent")
var applicationStartedEventId = context.GetEventId("github.com.genose.ApplicationStartedEvent")
var applicationReadyEventId = context.GetEventId("github.com.genose.ApplicationReadyEvent")
var applicationFailedEventId = context.GetEventId("github.com.genose.ApplicationFailedEvent")

func ApplicationEventId() context.ApplicationEventId {
	return applicationEventId
}

func ApplicationStartingEventId() context.ApplicationEventId {
	return applicationStartingEventId
}

func ApplicationEnvironmentPreparedEventId() context.ApplicationEventId {
	return applicationEnvironmentPreparedEventId
}

func ApplicationContextInitializedEventId() context.ApplicationEventId {
	return applicationContextInitializedEventId
}

func ApplicationPreparedEventId() context.ApplicationEventId {
	return applicationPreparedEventId
}

func ApplicationStartedEventId() context.ApplicationEventId {
	return applicationStartedEventId
}

func ApplicationReadyEventId() context.ApplicationEventId {
	return applicationReadyEventId
}

func ApplicationFailedEventId() context.ApplicationEventId {
	return applicationFailedEventId
}

type GenoseApplicationEvent interface {
	context.ApplicationEvent
	GetGenoseApplication() *GenoseApplication
	GetArgs() ApplicationArguments
}

type ApplicationStartingEvent struct {
	app       *GenoseApplication
	args      ApplicationArguments
	timestamp int64
}

func NewApplicationStarting(app *GenoseApplication, args ApplicationArguments) ApplicationStartingEvent {
	return ApplicationStartingEvent{
		app,
		args,
		time.Now().Unix(),
	}
}

func (event ApplicationStartingEvent) GetEventId() context.ApplicationEventId {
	return applicationStartingEventId
}

func (event ApplicationStartingEvent) GetParentEventId() context.ApplicationEventId {
	return applicationEventId
}

func (event ApplicationStartingEvent) GetSource() interface{} {
	return event.app
}

func (event ApplicationStartingEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationStartingEvent) GetGenoseApplication() *GenoseApplication {
	return event.app
}

func (event ApplicationStartingEvent) GetArgs() ApplicationArguments {
	return event.args
}

type ApplicationEnvironmentPreparedEvent struct {
	app         *GenoseApplication
	args        ApplicationArguments
	timestamp   int64
	environment core.ConfigurableEnvironment
}

func NewApplicationEnvironmentPreparedEvent(app *GenoseApplication, args ApplicationArguments, env core.ConfigurableEnvironment) ApplicationEnvironmentPreparedEvent {
	return ApplicationEnvironmentPreparedEvent{
		app,
		args,
		time.Now().Unix(),
		env,
	}
}

func (event ApplicationEnvironmentPreparedEvent) GetEventId() context.ApplicationEventId {
	return applicationEnvironmentPreparedEventId
}

func (event ApplicationEnvironmentPreparedEvent) GetParentEventId() context.ApplicationEventId {
	return applicationEventId
}

func (event ApplicationEnvironmentPreparedEvent) GetSource() interface{} {
	return event.app
}

func (event ApplicationEnvironmentPreparedEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationEnvironmentPreparedEvent) GetGenoseApplication() *GenoseApplication {
	return event.app
}

func (event ApplicationEnvironmentPreparedEvent) GetArgs() ApplicationArguments {
	return event.args
}

func (event ApplicationEnvironmentPreparedEvent) GetEnvironment() core.ConfigurableEnvironment {
	return event.environment
}

type ApplicationContextInitializedEvent struct {
	app       *GenoseApplication
	args      ApplicationArguments
	timestamp int64
	context   context.ConfigurableApplicationContext
}

func NewApplicationContextInitializedEvent(app *GenoseApplication, args ApplicationArguments, ctx context.ConfigurableApplicationContext) ApplicationContextInitializedEvent {
	return ApplicationContextInitializedEvent{
		app,
		args,
		time.Now().Unix(),
		ctx,
	}
}

func (event ApplicationContextInitializedEvent) GetEventId() context.ApplicationEventId {
	return applicationContextInitializedEventId
}

func (event ApplicationContextInitializedEvent) GetParentEventId() context.ApplicationEventId {
	return applicationEventId
}

func (event ApplicationContextInitializedEvent) GetSource() interface{} {
	return event.app
}

func (event ApplicationContextInitializedEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationContextInitializedEvent) GetGenoseApplication() *GenoseApplication {
	return event.app
}

func (event ApplicationContextInitializedEvent) GetArgs() ApplicationArguments {
	return event.args
}

func (event ApplicationContextInitializedEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return event.context
}

type ApplicationPreparedEvent struct {
	app       *GenoseApplication
	args      ApplicationArguments
	timestamp int64
	context   context.ConfigurableApplicationContext
}

func NewApplicationPreparedEvent(app *GenoseApplication, args ApplicationArguments, ctx context.ConfigurableApplicationContext) ApplicationPreparedEvent {
	return ApplicationPreparedEvent{
		app,
		args,
		time.Now().Unix(),
		ctx,
	}
}

func (event ApplicationPreparedEvent) GetEventId() context.ApplicationEventId {
	return applicationPreparedEventId
}

func (event ApplicationPreparedEvent) GetParentEventId() context.ApplicationEventId {
	return applicationEventId
}

func (event ApplicationPreparedEvent) GetSource() interface{} {
	return event.app
}

func (event ApplicationPreparedEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationPreparedEvent) GetGenoseApplication() *GenoseApplication {
	return event.app
}

func (event ApplicationPreparedEvent) GetArgs() ApplicationArguments {
	return event.args
}

func (event ApplicationPreparedEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return event.context
}

type ApplicationStartedEvent struct {
	app       *GenoseApplication
	args      ApplicationArguments
	timestamp int64
	context   context.ConfigurableApplicationContext
}

func NewApplicationStartedEvent(app *GenoseApplication, args ApplicationArguments, ctx context.ConfigurableApplicationContext) ApplicationStartedEvent {
	return ApplicationStartedEvent{
		app,
		args,
		time.Now().Unix(),
		ctx,
	}
}

func (event ApplicationStartedEvent) GetEventId() context.ApplicationEventId {
	return applicationStartedEventId
}

func (event ApplicationStartedEvent) GetParentEventId() context.ApplicationEventId {
	return applicationEventId
}

func (event ApplicationStartedEvent) GetSource() interface{} {
	return event.app
}

func (event ApplicationStartedEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationStartedEvent) GetGenoseApplication() *GenoseApplication {
	return event.app
}

func (event ApplicationStartedEvent) GetArgs() ApplicationArguments {
	return event.args
}

func (event ApplicationStartedEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return event.context
}

type ApplicationReadyEvent struct {
	app       *GenoseApplication
	args      ApplicationArguments
	timestamp int64
	context   context.ConfigurableApplicationContext
}

func NewApplicationReadyEvent(app *GenoseApplication, args ApplicationArguments, ctx context.ConfigurableApplicationContext) ApplicationReadyEvent {
	return ApplicationReadyEvent{
		app,
		args,
		time.Now().Unix(),
		ctx,
	}
}

func (event ApplicationReadyEvent) GetEventId() context.ApplicationEventId {
	return applicationReadyEventId
}

func (event ApplicationReadyEvent) GetParentEventId() context.ApplicationEventId {
	return applicationEventId
}

func (event ApplicationReadyEvent) GetSource() interface{} {
	return event.app
}

func (event ApplicationReadyEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationReadyEvent) GetGenoseApplication() *GenoseApplication {
	return event.app
}

func (event ApplicationReadyEvent) GetArgs() ApplicationArguments {
	return event.args
}

func (event ApplicationReadyEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return event.context
}

type ApplicationFailedEvent struct {
	app       *GenoseApplication
	args      ApplicationArguments
	timestamp int64
	context   context.ConfigurableApplicationContext
	err       error
}

func NewApplicationFailedEvent(app *GenoseApplication, args ApplicationArguments, ctx context.ConfigurableApplicationContext, err error) ApplicationFailedEvent {
	return ApplicationFailedEvent{
		app,
		args,
		time.Now().Unix(),
		ctx,
		err,
	}
}

func (event ApplicationFailedEvent) GetEventId() context.ApplicationEventId {
	return applicationFailedEventId
}

func (event ApplicationFailedEvent) GetParentEventId() context.ApplicationEventId {
	return applicationEventId
}

func (event ApplicationFailedEvent) GetSource() interface{} {
	return event.app
}

func (event ApplicationFailedEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationFailedEvent) GetGenoseApplication() *GenoseApplication {
	return event.app
}

func (event ApplicationFailedEvent) GetArgs() ApplicationArguments {
	return event.args
}

func (event ApplicationFailedEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return event.context
}

func (event ApplicationFailedEvent) GetError() error {
	return event.err
}
