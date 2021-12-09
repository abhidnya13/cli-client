# Demo CLI Client

This is a sample cli which can be used to do an interactive login with your Microsoft work/school accounts or personal microsoft accounts. It queries the graph endpoint for your basic information which is then dispalyed on the CLI.

To run the sample you can simply do:
``` Go
go run ./ login 
```

This will direct you to sign in and consent page for authentication. Once authentication completes, you can go back to the CLI to look at your information pulled from the graph. 
