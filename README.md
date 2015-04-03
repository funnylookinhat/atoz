# Atoz - Not A to Z.

**See: [http://en.memory-alpha.org/wiki/Atoz](http://en.memory-alpha.org/wiki/Atoz)**

##Atoz is a documentation generator for modern APIs.

Most APIs that are written today communicate exclusively in JSON objects - both requests and responses.  Writing documentation for these APIs within your code, however, can be a total pain.  Atoz makes this easier.

Atoz takes inline documentation that you add to your code and converts it into a giant, behemoth of a JSON object that you can parse into static HTML ( or even render on the fly if you're feeling gutsy ).  This output contains all of your "Actions" ( i.e. end-points ) and "Objects" ( i.e. models or types ).

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
 * @return {Integer} user.id The unique ID that represents this user in the system.
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
      "uri": "",
      "description": "Fetch a user from the application.",
      "parameters": [
        {
          "type": "integer",
          "limit": -1,
          "description": "The user id to lookup.",
          "children": [
            
          ]
        }
      ],
      "returns": [
        {
          "type": "object",
          "limit": -1,
          "description": "An object representing the user.",
          "children": [
            {
              "type": "string",
              "limit": 0,
              "description": "The user's email address.",
              "children": [
                
              ]
            },
            {
              "type": "integer",
              "limit": -1,
              "description": "The unique ID that represents this user in the system.",
              "children": [
                
              ]
            },
            {
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

