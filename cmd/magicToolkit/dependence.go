package main

func init() {
	actionList = append(actionList, &action{
		name:         "updateDep",
		defaultValue: false,
		description:  "update project dependence",
		funcPtr:      updateDep,
	})
}

func updateDep(curPath string) {

}
