<!--
  - A Golang Operator - Over Easy
    - Overview
    - Lab Specification (BDD Style)
    - Step-by-Step Detailed Lab Walkthrough
      - Summary
      - Scaffolding
      - Writing Your Operator Specifications and Status
      - Writing Your Operator Controller Implementation
      - Unit Testing
      - End to End Testing
-->
# A Golang Operator - Over Easy

## Environment Setup

Ensure lab pre-requisites have been met. See: [Lab Requirements](../01/03-lab-requirements.md)

## Lab Specifications (BDD Style)

### Story

- **DESCRIPTION**: An Operator with a single busy box pod that shuts down after a user specified amount of time
- **GIVEN**: A scaffolded operator
- **WHEN**: the specification `timeout` is added as an attribute to the operator
- **AND**: the specification `timeout` is set to a numeric value in seconds
- **AND**: and an Operator instance is created
- **THEN**: the busy box pod will remain available for the specified `timeout` in seconds,
- **AND**: log the message, `busybox pod expired` upon `timeout` expiration 
- **AND**: shutdown after the specified amount `timeout` duration

### Acceptance Criteria

- [ ] - the Operator must have a `timeout` specification attribute
- [ ] - the Operator instance must start a busybox pod for the duration of `timeout` in seconds
- [ ] - the Operator instance must shut down after duration of `timeout` in seconds, has expired
- [ ] - the Operator instance must log a message before expiration

## Lab Walkthrough

