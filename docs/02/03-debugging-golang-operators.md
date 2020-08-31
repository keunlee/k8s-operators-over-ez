# Development Environments

The following are prescribed options that you may choose to follow. There is no hard requirement to setup your environment as prescribed.

# Debugging GoLang Operators with Visual Studio Code

1. Open this code at the root of the lab code directory (`docs/02/labs/000/code`) with Visual Studio Code. 
2. Install necessary VS Code plugins/extensions for Go:  https://github.com/golang/vscode-go#getting-started
3. Follow directions for debugging: https://github.com/golang/vscode-go/blob/master/docs/debugging.md#set-up


To debug Unit tests: 
- Set break points
- Select the file `suite_test.go` as your "Run and Debug" starting point

To debug an operator instance:
- Set break points 
- Select the file `main.go` as your "Run and Debug" starting point
- Create an instance of an operator: ` k apply -f config/samples/operators-over-ez_v1alpha1_opsovereasy.yaml`

For this setup, you will not need a `launch.json` file.

# Debugging GoLang Operators with JetBrains Intellij/GoLand

1. Install Go plugins
2. Create new project project and specify: 
- Go
- Location: Root of the lab code directory (`docs/02/labs/000/code`) 
- GOROOT: Specify or install one if necessary (i.e. Go 1.14.6)
- Check `Index entire GOPATH`
- Click `Create`
- Select `Create from existing sources`

To debug Unit tests: 
- Set break points
- Right-click the file `suite_test.go` and select: `Debug 'suite_test.go'`

To debug an operator instance:
- Set break points 
- Right-click the file `main.go` and select: `Debug 'main.go'`
- Create an instance of an operator: ` k apply -f config/samples/operators-over-ez_v1alpha1_opsovereasy.yaml`


[Return to Table of Contents](../../../../)

