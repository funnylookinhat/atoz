# Atoz - Not A to Z.

**See: [http://en.memory-alpha.org/wiki/Atoz](http://en.memory-alpha.org/wiki/Atoz)**

Atoz is the result of my looking for a documentation tool that reflects the way 
I believe REST APIs should be designed: JSON in and JSON out.  More specifically,
your code should work with your application in nearly the exact same way that 
a third party client might.  I wanted a simple way to document these sorts of 
interactions without adding any significant amount of weight to my own code.

##Atoz is a documentation generator for modern APIs.

Most APIs that are written today communicate exclusively in JSON objects - both 
requests and responses.  Writing documentation for these APIs within your code, 
however, can be a total pain.  Atoz makes this easier.

Atoz takes inline documentation that you add to your code and converts it into 
a giant, behemoth of a JSON object that you can parse into static HTML ( or 
even render on the fly if you're feeling gutsy ).  This output contains all of 
your "Actions" ( i.e. end-points ) and "Objects" ( i.e. models or types ).

Here's an example:

```
<?php 

/**
 * ---ATOZAPI---
 * @name Get User
 * @ref /MyApp/User/Get
 * @uri /api/user/get
 * @description Fetch a user from the application.
 * @required {Integer} id The user id to lookup.
 * @return {Object} user An object representing the user.
 * @return {Integer} user.id
 * @return {String} user.name The user's name.
 * @return {String,254} user.email The user's email address.
 * ---ATOZEND---
 */

function getUser($id) {
    // Do something awesome here and return a user.
}
```

Running Atoz against that source tree would produce the following:

```
{
  "actions": [
    {
      "name": "Get User",
      "ref": "\/MyApp\/User\/Get",
      "uri": "\/api\/user\/get",
      "description": "Fetch a user from the application.",
      "parameters": [
        {
          "name": "id",
          "flag": "required",
          "type": "integer",
          "limit": -1,
          "description": "The user id to lookup.",
          "children": [
            
          ]
        }
      ],
      "returns": [
        {
          "name": "user",
          "flag": "",
          "type": "object",
          "limit": -1,
          "description": "An object representing the user.",
          "children": [
            {
              "name": "email",
              "flag": "",
              "type": "string",
              "limit": 254,
              "description": "The user's email address.",
              "children": [
                
              ]
            },
            {
              "name": "id",
              "flag": "",
              "type": "integer",
              "limit": -1,
              "description": "",
              "children": [
                
              ]
            },
            {
              "name": "name",
              "flag": "",
              "type": "string",
              "limit": 0,
              "description": "The user's name.",
              "children": [
                
              ]
            }
          ]
        }
      ]
    }
  ],
  "objects": [
    
  ]
}
```

## Usage

By default, Atoz will output JSON to stdout and will recurse the current 
working directory.  Passing `-dir some/path` will search the provided directory 
instead of the current one.  Additionally, you can specify 
`-output some/file.json` to write the JSON directly to a file.

`./atoz -dir path/to/source/tree -output some/json/file.json`

Atoz will recursively search through the provided directory for valid UTF-8 
encoded text files that include definitions, actions, or objects.  These all 
start with a line that includes one of the following: `---ATOZAPI---`, 
`---ATOZOBJ---`, or `---ATOZDEF---`.  Any group of lines must be terminated 
with a line containing `---ATOZEND---`.

So - as an example - an object definition might look like this:

```
/**
 * ---ATOZOBJ---
 * @name Some name
 * @ref Some reference
 * ---ATOZEND---
 */
```

Atoz supports three types of definitions: Actions, Objects, and Definitions.  
Ignore the poor naming convention of the last type.

## Actions

Actions represent API end-points - they don't specify whether they're GET, POST, or PUT, 
as Atoz assumes that every request body includes a JSON object, and every response 
body is also a JSON object.  In most cases, people will implement this with only 
POST requests - but semantically changing the HTTP Request to PUT for Updates or 
something of that sort can certainly be added at a later point in time.

Actions support the following attributes:
- `@name Value` A title for the action.
- `@ref Value` A canonical reference.
- `@uri Value` The expected URI for the api-end point over HTTP.
- `@description Value` 
- `@parameter {Type,Limit} Object.Space Description` A parameter that can be sent to the action.
- `@required {Type,Limit} Object.Space Description` A required parameter.
- `@optional {Type,Limit} Object.Space Description` An optional parameter.
- `@returns {Type,Limit} Object.Space Description` A value that is returned.
- `@success {Type,Limit} Object.Space Description` A value returned only upon success.
- `@failure {Type,Limit} Object.Space Description` A value returned only on failure.

## Objects

Objects are the models that your API frequently works with.  These are defined 
similarly to Actions, but they offer fewer fields.

- `@name Value` A title for the object.
- `@ref Value` A canonical reference.
- `@description Value`
- `@property {Type,Limit} Object.Space Description` Any key/value stored in this object.

## Definitions

Definitions are subsets of lines that you can define once and include anywhere 
else.  These allow you to avoid repeating common requirements for API end-points 
or frequent parameters for objects.  They can have any of the defined lines 
that either objects or actions have.

For example, a common set of authorization parameters for your API Actions might 
be defined as follows:

```
/**
 * ---ATOZDEF---
 * @ref /Action/AuthParams
 * @required {Object} auth
 * @required {Integer} auth.id User's ID.
 * @required {String} auth.token Unique access token for this user session.
 * ---ATOZEND---
 */
```

If you wanted to add this to an action, you simply add an `@include` statement 
with it's reference.

```
/**
 * ---ATOZAPI---
 * @name User Lookup
 * @ref /User/Lookup
 * @uri /api/user/lookup
 * @description Fetch information on a user based on ID.
 * @include /Action/AuthParams
 * @required {Integer} id The ID of the user you are requesting.
 * ---ATOZEND---
 */
```

In reality, this pulls that definition's lines in (minus the `@ref` statement) 
in place of the `@include` statement.  You should assume that any value that is 
pulled in will overwrite a previous value if they use the same object.space or 
reserved key.

## String Values

There are two types of lines that are used when specifying definitions in Atoz. 
For the following, Atoz expects everything following the line type to be a 
single string.

- `@name` 
- `@ref` 
- `@uri` 
- `@description` 

Any of the following are perfectly valid:

- `@name Some name for a function`
- `@ref /A/Canonical/Path`
- `@ref A String to Reference` - Note that @ref is used to identify other 
objects, actions, and definitions as a key - so including spaces, while 
functional, might not be the best decision.
- `@description A **markdown** string of text to ~~print~~ parse later.`

Not that Atoz will not do any further parsing of text (i.e. support for 
markdown), but you could simply parse it further when generating HTML or 
whatever means of documentation you wish.

## Parameters, Returns, and Properties

For all of the other line types, Atoz expects you to define a key for a value. 
These are formatted generally as `@linetype {Type,Limit} Object.space Description`:

**Type** can be any of the following:
- `Boolean`
- `Integer`
- `Decimal`
- `String`
- `Array`
- `Object`

Remember, these are tools meant to help identify expected input/output for 
your API.  As an example, there are obviously many more types of numbers that 
could be more accurately specified (Float, Double, Integer, Unsigned, etc.), 
but the goal is to tell the user what type of data to send, not how it will be 
handled or how they should handle it.

The **Limit** parameter is used only for a few of these line types.  If one is 
not supported for a specific Type, the resulting JSON will have `-1` for it's 
value. If one is supported and not provided, the assumption is `0`, which should 
be translated as "unlimited".

The following Types have limits:

- `Decimal` - Number of decimal points in precision that can be provided.
- `String` - Maximum length of a string.
- `Array` - Maximum number of elements that can be in the array.

**Object.space** is used to show where in the JSON object a key should be placed. 
For example, if I had an integer at the root of my return object, I might specify 
it like this:

`@returns {Integer} id`

This means that I could expect the resulting JSON to produce something along the 
lines of this:

```
{
	"id": 123456,
	...
}
```

Furthermore, these can be used to show nested key/values:

```
@returns {Object} user
@returns {Integer} user.id
@returns {String} user.name
```

This would indicate an expected result similar to this:

```
{
	"user": {
		"id": 123456,
		"name": "John Doe"
	},
	...
}
```

A couple important things to note:

- **Object** types must be declared for their children to be parsed.  You 
cannot specify `user.id` without specifying `user` as an Object.  They don't have 
to be in any particular order, but it would help semantically for anyone trying 
to manage your source code.
- **Array** types act exactly like **Object** types, except that you're specifying 
a list of objects instead of a single one.  That is to say, you can only have 
an array of Objects.
- **Description** values are optional.  In many cases, adding anything to the 
actual name of the key / value is excessive ( i.e. describing what user.id 
might mean).  If no description is provided, the resulting JSON will just have 
a blank string.

Lastly, parameters and return values can be described in a few different ways.
By default, you can simply used `@parameter` and `@returns` to generally show 
what information is coming in and going out of your API.  However, many APIs will 
have specific information only upon success or a failure, or might require some 
parameter while another is optional.  You can describe these as follows:

- `@parameter` can also be:
    - `@required` - A value that is required by the API end-point.
    - `@optional` - One that is specifically optional.
- `@returns` can also be:
    - `@success` - A value returned only upon success.
    - `@failure` - A value returned only if the request fails.

These will still show up in the same place you would expect `@parameter` and 
`@returns` values within the resulting JSON, but they will have a value 
for `flag` specified.  For example:

```
@required {String,254} email
```

Would produce:

```
{
	...
	{
		"name": "email",
		"flag": "required",
		"type": "string",
		"limit": 254,
		"description": "",
		"children": []
	},
	...
}
```

You can use these flags to apply special classes to whatever HTML you might 
generate to help users identify those unique points of your API.