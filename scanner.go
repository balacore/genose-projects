package genose

import (
	"fmt"
	"github.com/genose-projects/goo"
	context "github.com/genose-projects/genose-context"
	core "github.com/genose-projects/genose-core"
	peas "github.com/genose-projects/genose-peas"
)

type componentScanner struct {
}

func newComponentScanner() componentScanner {
	return componentScanner{}
}

func (scanner componentScanner) scan(contextId context.ContextId, logger context.Logger) (int, error) {
	processors, err := scanner.getProcessorInstances()
	if err != nil {
		return -1, nil
	}
	var componentCount = 0
	result := core.ForEachComponentType(func(componentName string, componentType goo.Type) error {
		logger.Trace(contextId, fmt.Sprintf("Registered component %s", componentName))
		err := scanner.checkComponent(componentType, processors)
		if err != nil {
			return err
		}
		componentCount++
		return nil
	})
	return componentCount, result
}

func (scanner componentScanner) checkComponent(componentType goo.Type, processors []interface{}) (err error) {
	for _, processorInstance := range processors {
		if processor, ok := processorInstance.(core.ComponentProcessor); ok {
			if processor.SupportsComponent(componentType) {
				err = processor.ProcessComponent(componentType)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (scanner componentScanner) getProcessorInstances() ([]interface{}, error) {
	var instances []interface{}
	result := core.ForEachComponentProcessor(func(processorName string, processorType goo.Type) error {
		instance, err := peas.CreateInstance(processorType, []interface{}{})
		if err != nil {
			return err
		}
		instances = append(instances, instance)
		return nil
	})
	return instances, result
}
