## Squirrelistic JSON Preprocessor ##

JSON files can be long and boring, overgrowth of lines, slowly hypnotising, with its repetitive cadence and sea of malevolent characters.<br/>
That's where DJ Squirrel shows up - to alleviate the JSON pain with it's flashy pre-processing formula:

**Squirrelistic JSON Template (SJT) + Input Parameters (JSON) = Long JSON**

This equation probably doesn't make sense, so let's delve into an obligatory Hello World example.

*Template:*

```JSON
{
    "hello": "{{ .helloText }}"
}
```

*Plus parameters:*

```JSON
{
    "helloText": "world"
}
```

*Equals:*

```JSON
{
    "hello": "world"
}
```

## Examples ##

Check out *[test-files](test-files)* directory<br/>
&nbsp;&nbsp;in this Github repository<br/>
Can you see the examples factory?<br/>
&nbsp;&nbsp;is it satisfactory?

## Command line syntax ##

```Shell
sjcp -t=template.sjt -p=params.json -o=result.json
```

## Squirrelistic JSON Template features ##

- Conditional logic (if, else if, else)
- Includes, including includes including includes (iiii)
- Loops
- Comments
- And much more...

## Usage scenarios ##

- Reusable JSON snippets
- Splitting large JSON files into small logical units
- More readable Azure ARM Templates
- Interspecies communication with Alien Squirrels

## Engine ##

Template engine piggybacks on [Golang Text Template package](https://golang.org/pkg/text/template/) so all its syntax is valid as long as it produces valid JSON file.

The following functions were added on top:
- {{ include filename }} - includes another SJT template (SJT buddy), which in turn can include more buddies
- {{ json .paramName }} - JSON representation of parameter (see the explanation bellow)
- {{ jsonEscape .paramName }} - escapes JSON offending characters (double quote, newline etc.) lest it ends up in JSON prison.

### The {{ .param }} vs {{ json .param }} vs {{ jsonEscape .param }} conundrum ###

Let's have this nifty parameters file where:
- Param1 is string which contains single double quote, which is encoded as \" otherwise JSON gods would be angry.
- Param2 is an array. A nice one.

```JSON
{
    "param1": "\"",
    "param2": [ "Squirrels", "are", "nice" ]
}
```

Here is the sexy table which will shed a light on the JSON encoding multiverse.

Template snippet | Output | Explanation
---------------- | ------ | -----------
{{&nbsp;.param1&nbsp;}} | " | Not escaped, not enclosed in ""
{{&nbsp;json&nbsp;.param1&nbsp;}} | "\\"" | Escaped, enclosed in ""
{{&nbsp;jsonEscape&nbsp;.param1&nbsp;}} | \\" | Escaped, not enclosed in ""
{{&nbsp;.param2&nbsp;}} | [Squirrels are nice] | String representation of an array, not escaped, not enclosed in ""
{{&nbsp;json&nbsp;.param2&nbsp;}} | \["Squirrels","are","nice"\] | JSON representation of an array, escaped and everything
{{&nbsp;jsonEscape&nbsp;.param2&nbsp;}} | **Kernel Panic!** | Invalid, as jsonEscape only accepts strings, not array or other objects

## Contributions ##

All contributions, suggestions, bug reports and peanuts are welcome.<br/>
Bear in mind that I wrote this project in 2 days, after learning Golang from tutorial, so The Perfection is not the name of this game.

## Links ##

 - [Golang Text Template](https://golang.org/pkg/text/template/)
 - [Golang Templates Cheatsheet](https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet)