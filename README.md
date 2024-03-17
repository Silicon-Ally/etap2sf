# eTapestry Migration Tool

This is a tool for migrating data from eTapestry to Salesforce.
It is comprehensive and highly customizable.

Though this repository solves 90-95% of the migration for an average
organization, getting the nuances of any migration correct requires someone with
an understanding of the business use cases, and some basic proficiency in the
Go programming language - please don't pursue a migration using this repository
without someone who can help you with each of these challenges.

Do you need someone to be that helping hand? Reach out to us at 
`etap-migration@siliconally.org`. We're a small software-nonprofit, and we've
done this migration before.

## What does a migration look like?

Broadly, the migration tookit is broken into a three major phases, with 20 discrete steps in total, each of which is idempotent (safe to repeat).

1. [Steps 1-3] Export data from eTapestry, using a combination of the [SOAP API](https://app.etapestry.com/hosted/files/api3/home.html) and raw HTTP requests (for attachments).
  * This process can take quite a bit of time if your data is large (many hours is typical).
2. [Steps 4-18] Create a [Salesforce Sandbox](https://www.salesforce.com/products/sandboxes-environments/), and upload a sample of your data into that sandbox.
  * Because Salesforce sandboxes are disposable, this process can be done a large number of times, allowing you to test various customizations and configurations in a safe and iterative way.
3. [Steps 4-20] Uploading the final version to Salesforce - this is largely the same as phase #2, with a few nuances because of data protections on production Salesforce instances.

## How to get started?

This codebase is composed of a number of steps - each step is a Go binary
that you build and execute. You can see the full set of steps in the `steps`
directory. When a step fails because it needs some action from you, it'll
typically give you instructions on what you need to do.

Start off by authenticating to Salesforce and eTapestry by running step 1:

```
go run steps/step01_validate_api_access
```

## Questions, Concerns, Suggestions, Bugs, etc.

For any and all commentary, support, or PRs, please communicate through GH Issues.
