package loggernew

import (
	"log"

	"github.com/google/uuid"
)

func (g *CrashGuardClient) new(config *LoggingConfig) (err error) {

	g.Config = config
	return nil
}

func (g *CrashGuardClient) log(object *InformationConstruct, severity string, logTag string) {

	defer func(object *InformationConstruct, severity string, logTag string) {
		if r := recover(); r != nil {
			if object.Operation != nil {
				log.Println("CRASHGUARD LOGGER FAILED, OP ID:", object.Operation.ID, "\n", r)
			} else {
				object.Operation = &Operation{ID: uuid.New().String()}
				log.Println("CRASHGUARD LOGGER FAILED, OP ID:", object.Operation.ID, "\n", r)
			}
			object.print(logTag, severity, g.Config.Debug)
		}
	}(object, severity, logTag)

	stacktrace, err := getStack(g.Config)
	if err != nil {
		log.Println(err) // handle this better
	}
	if stacktrace != "" {
		object.StackTrace = stacktrace
	}

	object.print(logTag, severity, g.Config.Debug)

}

func (g *CrashGuardClient) close() {
	// no op
}
