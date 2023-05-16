# tdd-ddd-go
tdd code from ddd with go

Step
- Create a simple failed test
    - without input
    - just output
- Implement simple code to pass the test
- improve test, add the input. run test failed
- implement code to pass the test
- in test, add mock dependency used in the code
    - add input, output expectation, and how many it calls when using mock
        - we need to think ahead about the requirement
        - what is the input and output, and how many it calls for the dependency that we want to mock
    - run test -> failed
- in the code, add interface, and call the dependency
- if only success case, then we can just ignore the return value from dependency, including error
- if we want to test error value, then in the code we can add the error
