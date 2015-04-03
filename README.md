# Atoz - Not A to Z.

**See: [http://en.memory-alpha.org/wiki/Atoz](http://en.memory-alpha.org/wiki/Atoz)**

Atoz is the result of my looking for a documentation tool that reflects the way 
I believe REST APIs should be designed: JSON in and JSON out.  I wanted a simple 
way to generate these sorts of interactions without adding any significant 
amount of weight to my own code.

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
 * @description Fetch a user from the application.
 * @required {Integer} id The user id to lookup.
 * @return {Object} user An object representing the user.
 * @return {Integer} user.id
 * @return {String} user.name The user's name.
 * @return {String} user.email The user's email address.
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
              "limit": 0,
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

Atoz will recursively search through a provided directory for valid UTF-8 
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

## Definitions

## Parameters, Returns, and Properties

