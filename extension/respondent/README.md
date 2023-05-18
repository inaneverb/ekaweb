Respondent is an "error handler" middleware. A special middleware
that will be executed IMMEDIATELY AFTER some middleware or route's handler
returns an error. Respondent is a worker for that returned error.
It allows you to replace one error to another, convert an error to special struct
and then use that struct to generate HTTP response. Respondent is a powerful
mechanism to convert error to some message, an API user's will see.

First of all, the Respondent has 3 things: Replacer, Expander, Applicator.
All of them are interfaces, thus you can use their default implementations
(named with "Common" prefixed) or implement your own variants.
Only Expander is REQUIRED. The rest of them, Replacer and Applicator can be omitted.
If they're omitted their default variants will be used.
Keep in mind, that the flow is: Replacer -> Expander -> Applicator.

Expander is used to convert incoming error into Manifest - a special object,
that represents an error in a terms of HTTP response.
The main idea of CommonExpander is:
Bind error A with manifest B and when error A is occurred, return manifest B.
Thus all manifest are static objects that are created at the Respondent creation.
But Manifests may have "extractors" for "reason(s)" field(s). It makes a way
to fill these fields dynamically depends on each HTTP request.

Replacer is used to replace one error with another. And, if present, do it
BEFORE Expander will start their job. The CommonReplacer just does that thing,
allowing you to specify how errors should be compared, directly or by Golang 1.13+
errors' API.

At the end, the Applicator takes Manifest from Expander, finally converts it
to the HTTP response and sends that HTTP response to the API user.
That's what CommonApplicator do performing JSON encoding of Manifest.

Again. Responder is: Replacer, then Expander, then Applicator.
If Replacer is omitted, the default (empty) is used (making no replaces).
If Applicator is omitted, the CommonApplicator is used.