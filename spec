General reference:

Objects: #/Object/Reference#
Actions: !/Api/Reference!

Definition Types:
@ref value
@required {Type,Limit}


Parameter / Response Types:
BOOLEAN - No Limit
INTEGER - No Limit
DECIMAL - Limit = Decimal Precision
STRING - Limit = Max Length
ARRAY - Limit = Max Length
OBJECT - No Limit


/**
 * ---ATOZDEF---
 * @ref /Defs/Authorization
 * @required {OBJECT} auth 
 * @required {INTEGER} auth.id 
 * @required {STRING,64} auth.key 
 * ---ATOZEND---
 */

/**
 * ---ATOZDEF---
 * @name /Defs/BaseResult
 * @success {BOOLEAN} success A boolean to show whether or not the request was successful.
 * @error {STRING} error An error message describing what went wrong.
 */

/**
 * ---ATOZAPI---
 * @name Lookup
 * @ref /MyApp/User/Lookup
 * @uri /User/Lookup
 * @description Get the information for a user.
 * @include /Defs/Authorization
 * @required {INTEGER} id The ID of the user.
 * @include /Defs/BaseResult
 * @success {#/Application/User#} user
 * ---ATOZEND---
 */

/**
 * ---ATOZOBJ---
 * @name User
 * @namespace /Application
 * @description A user in the application.
 * @property id INTEGER Unique ID of the user.
 * @property name STRING Name of the user.
 * @property email STRING Email address for the user.
 * ---ATOZEND---
 */
