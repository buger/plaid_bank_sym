### Task

In order to retrieve financial data, Plaid works with thousands of financial
institutions and banks. While modern banks have evolved to serve a plethora of
functions, at their core, banks must provide certain basic features. Today,
your task is to build the basic HTTP API for one of those banks!

Imagine you are designing a backend API for bank employees. It could ultimately
be consumed by multiple frontends (web, iOS, Android etc). There should be API
routes that allow them to:

  - Create a new bank account for a customer, with an initial deposit amount. A
    single customer may have multiple bank accounts.
  - Transfer amounts between any two accounts, including those owned by
    different customers.
  - Retrieve balances for a given account.
  - Retrieve transfer history for a given account.

Feel free to pre-populate your customers with the following:

```json
[
  {
    "id": 1,
    "name": "Jane Woods"
  },
  {
    "id": 2,
    "name": "Michael Li"
  },
  {
    "id": 3,
    "name": "Heidi Hasselbach"
  },
  {
    "id": 4,
    "name": "Rahul Pour"
  }
]
```

You are expected to design any other required models and routes for your API.

Treat this like a production system. We will be judging you on:

  - Completeness: did you complete the features?
  - Correctness: does the functionality act in sensible, thought-out ways?
  - Maintainability: is it written in a clean, maintainable way?
  - Testing: is the system adequately tested?
  - Documentation: is the API well-documented?

### Timeline

We're expecting this to take between 4 and 6 hours to complete. You'll be given
a 12 hour window to submit the task.

### Submission

Please send us two links:

- Link to the GitHub repository with your code.
- Link to a deployed sample API. If you do not have a deployable app ready at
  the end of the allotted time, you may skip this step. Feel free to use the
  hosting provider of your choice.
