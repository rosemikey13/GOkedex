module github.com/rosemikey13/GOkedex

go 1.22.4

require (
    github.com/rosemikey13/GOkedex/internal/poke-api v0.0.0
    github.com/rosemikey13/GOkedex/internal/pokecache v0.0.0
)

replace github.com/rosemikey13/GOkedex/internal/poke-api => ./internal/pokeApi/
replace github.com/rosemikey13/GOkedex/internal/pokecache => ./internal/pokecache/
