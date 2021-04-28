# Reference
`reference` is a "microservice". It consists of two main  sections of code:
- the `client`, which exposes what the rest of the world needs to interact 
  with the `service`;
- the `service`, which contains the implementation of the inner logic.

Definition of requests will be found in the `client`. They're written using 
protobuf, and can be regenerated with a `make generate` at the root of the 
project.