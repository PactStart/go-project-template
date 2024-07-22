package application

type ApplicationContext struct {
	ComponentMap map[string]interface{}
}

var AppContext ApplicationContext

func NewApplicationContext() {
	AppContext = ApplicationContext{
		ComponentMap: make(map[string]interface{}),
	}
}

func (e ApplicationContext) RegisterComponent(name string, component interface{}) {
	e.ComponentMap[name] = component
}

func (e ApplicationContext) GetComponent(name string) interface{} {
	return e.ComponentMap[name]
}
