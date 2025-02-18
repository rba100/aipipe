# Pokedex CLI

A commandline utility written in Golang that tells you about pokemon.

## Usage

`pokedex {pokemonname}` => returns matching pokemon and its abilities.

## Implementation

Calls `https://pokeapi.co/api/v2/pokemon/{pokemonname}`

which should return json that looks like this:

```json
{
  "name" : "ditto",
  "order" : 214,
  "abilities": [
    {
      "ability": {
        "name": "limber",
      },
      "is_hidden": false,
      "slot": 1
    },
    {
      "ability": {
        "name": "imposter",
      }
    }
  ]
}
```
Note the nested .abilities[].ability.name

This should be parsed and printed to the console like:
{name}:
 -{abilities[0].name}
 -{abilities[n].name}