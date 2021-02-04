package genose

import (
	"errors"
	context "github.com/genose-projects/genose-context"
	web "github.com/genose-projects/genose-web"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testGenoseApplicationEvent(t *testing.T,
	event GenoseApplicationEvent,
	eventId context.ApplicationEventId,
	parentEventId context.ApplicationEventId,
	source interface{},
	application *GenoseApplication,
	args ApplicationArguments) {

	assert.Equal(t, eventId, event.GetEventId())
	assert.Equal(t, parentEventId, event.GetParentEventId())
	assert.Equal(t, source, event.GetSource())
	assert.NotEqual(t, int64(0), event.GetTimestamp())
	assert.Equal(t, application, event.GetGenoseApplication())
	assert.Equal(t, args, event.GetArgs())
}

func TestApplicationStartingEvent(t *testing.T) {
	var application = NewGenoseApplication()
	var appArgs = getApplicationArguments(nil)
	event := NewApplicationStarting(application, appArgs)

	testGenoseApplicationEvent(t, event, ApplicationStartingEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, appArgs, event.GetArgs())
}

func TestApplicationEnvironmentPreparedEvent(t *testing.T) {
	var application = NewGenoseApplication()
	var appArgs = getApplicationArguments(nil)
	var environment = web.NewStandardWebEnvironment()
	event := NewApplicationEnvironmentPreparedEvent(application, appArgs, environment)

	testGenoseApplicationEvent(t, event, ApplicationEnvironmentPreparedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, environment, event.GetEnvironment())
}

func TestApplicationContextInitializedEvent(t *testing.T) {
	var application = NewGenoseApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewGenoseServerApplicationContext("app-id", "context-id")
	event := NewApplicationContextInitializedEvent(application, appArgs, ctx)

	testGenoseApplicationEvent(t, event, ApplicationContextInitializedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
}

func TestApplicationPreparedEvent(t *testing.T) {
	var application = NewGenoseApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewGenoseServerApplicationContext("app-id", "context-id")
	event := NewApplicationPreparedEvent(application, appArgs, ctx)

	testGenoseApplicationEvent(t, event, ApplicationPreparedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
}

func TestApplicationStartedEvent(t *testing.T) {
	var application = NewGenoseApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewGenoseServerApplicationContext("app-id", "context-id")
	event := NewApplicationStartedEvent(application, appArgs, ctx)

	testGenoseApplicationEvent(t, event, ApplicationStartedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
}

func TestApplicationReadyEvent(t *testing.T) {
	var application = NewGenoseApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewGenoseServerApplicationContext("app-id", "context-id")
	event := NewApplicationReadyEvent(application, appArgs, ctx)

	testGenoseApplicationEvent(t, event, ApplicationReadyEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
}

func TestApplicationFailedEvent(t *testing.T) {
	var application = NewGenoseApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewGenoseServerApplicationContext("app-id", "context-id")
	var err = errors.New("test error")
	event := NewApplicationFailedEvent(application, appArgs, ctx, err)

	testGenoseApplicationEvent(t, event, ApplicationFailedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
	assert.Equal(t, err, event.GetError())
}
