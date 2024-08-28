module github.com/rosemikey13/pokedex

go 1.22.4

require (
    github.com/rosemikey13/pokedex/internal/poke-api v0.0.0
     github.com/rosemikey13/pokedex/internal/pokecache v0.0.0
)

replace github.com/rosemikey13/pokedex/internal/poke-api => ./internal/pokeApi/
replace github.com/rosemikey13/pokedex/internal/pokecache => ./internal/pokecache/
