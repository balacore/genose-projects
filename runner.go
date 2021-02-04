package genose

import context "github.com/genose-projects/genose-context"

type ApplicationRunner interface {
	OnApplicationRun(context context.Context, arguments ApplicationArguments)
}
