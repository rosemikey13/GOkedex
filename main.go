package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	pokeapi "github.com/rosemikey13/GOkedex/internal/poke-api"
	pokecache "github.com/rosemikey13/GOkedex/internal/pokecache"
)

type command struct {
	name string
	description string
	callback func(*pokeapi.Config, string) error
}

var cache = pokecache.NewCache(10)
var pokedex = map[string]pokeapi.Pokemon{}

func main() {

	config := pokeapi.New()
	
	
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		choice := strings.Split(scanner.Text(), " ")[0]
		name := ""

		if choice == "explore" || choice == "catch" || choice == "inspect" {
		name = strings.Split(scanner.Text(), " ")[1]
		}

		if _, ok := commands[choice]; !ok{
			fmt.Printf("\n \"%v \" is not a valid choice, Please try again.\n\n", choice)
			continue
		}
		err := commands[choice].callback(config, name)
		if err != nil {
			fmt.Printf("\nError: %v", err)
		}
	}
}

func commandHelp(p *pokeapi.Config, name string) error {
	commands := getCommands()
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Printf("\n%v: %v\n", commands["help"].name, commands["help"].description)
	fmt.Printf("\n%v: %v\n", commands["map"].name, commands["map"].description)
	fmt.Printf("\n%v: %v\n", commands["mapb"].name, commands["mapb"].description)
	fmt.Printf("\n%v: %v\n", commands["explore"].name, commands["explore"].description)
	fmt.Printf("\n%v: %v\n\n", commands["exit"].name, commands["exit"].description)
	
	return nil
}

func commandExit(p *pokeapi.Config, name string) error {
	os.Exit(0)
	return nil
}

func commandMap(p *pokeapi.Config, name string) error{
	if p.Next == nil{
		return errors.New("There is no page after this one, use the \"mapb\" command to view previous pages.\n")
	}
	
	
	if data, found := cache.Get(*p.Next); found{
		locations, err := pokeapi.ParseLocations(data)
		p.Previous = locations.Previous
		p.Next = locations.Next
		
			if err != nil {
				fmt.Printf("Error getting locations: %v", err)
				os.Exit(1)
			}
		
		
			for _, location := range locations.Results{
				fmt.Println(location.Name)
			}
		return nil
	}
	
		
	
	
	data, err := pokeapi.GetLocations(*p.Next)
	locations, err := pokeapi.ParseLocations(data)
	cache.Add(*p.Next, data)
	
	p.Previous = locations.Previous
	p.Next = locations.Next


	for _, location := range locations.Results{
		fmt.Println(location.Name)
	}

	return err
}


func commandMapB(p *pokeapi.Config, name string) error{
	if p.Previous == nil{
		return errors.New("There is no page before this one, use the \"map\" command to view the proceeding pages.\n")
	}

	if data, found := cache.Get(*p.Previous); found{
		locations, err := pokeapi.ParseLocations(data)
		p.Previous = locations.Previous
		p.Next = locations.Next

		if err != nil {
			fmt.Printf("Error getting locations: %v", err)
			os.Exit(1)
		}

		for _, location := range locations.Results{
			fmt.Println(location.Name)
		}
		
		return nil
	}
	
	data, err := pokeapi.GetLocations(*p.Previous)
	locations, err := pokeapi.ParseLocations(data)
	cache.Add(*p.Previous, data)

	p.Previous = locations.Previous
	p.Next = locations.Next

	for _, location := range locations.Results{
		fmt.Println(location.Name)
	}

	return err
}

func commandExplore(p *pokeapi.Config, name string) error {

	fmt.Printf("\nExploring %v...\n", name)
	
	if data, found := cache.Get("https://pokeapi.co/api/v2/location-area/" + name); found{
		encounters, err := pokeapi.ParseAreaPokemon(data)
		
		if err != nil {
			fmt.Printf("Error getting pokemon: %v", err)
			os.Exit(1)
		}

		fmt.Println("Found Pokemon:")

		for _, encounter := range encounters.PokemonEncounters{
			fmt.Println(" - " + encounter.Pokemon.Name)
		}
		
		return nil
	}
	
	data, webAddress, err := pokeapi.GetAreaPokemon(name)

	if err != nil {
		return fmt.Errorf("Error retrieving Pokemon Data: %v", err)
	}

	cache.Add(webAddress, data)

	encounters, err := pokeapi.ParseAreaPokemon(data)
		

		if err != nil {
			fmt.Printf("Error getting pokemon: %v", err)
			os.Exit(1)
		}

		for _, encounter := range encounters.PokemonEncounters{
			fmt.Println(" - " + encounter.Pokemon.Name)
		}
	
		return nil
}

func commandCatch(p *pokeapi.Config, name string) error {
	fmt.Printf("Throwing a Pokeball at %v...", name)
	pokemon, err := pokeapi.GetPokemonInfo(name)
	if err != nil {
		return fmt.Errorf("\nUnable to retrieve info for pokemon named: %v.\nAre you sure this Pokemon exists?\n", name)
	}

	catchAttempt := rand.Int63n(350)
	if catchAttempt >= int64(pokemon.BaseExperience){
		fmt.Printf("\n%v was caught!\n", name)
		pokedex[name] = pokemon
		return nil
	}

	fmt.Printf("\n%v escaped!\n", name)


	return nil
}

func commandInspect(p *pokeapi.Config, name string) error {
	
	if pokemon, exists := pokedex[name]; exists {
		fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\nStats:\n", pokemon.Name, pokemon.Height, pokemon.Weight)
		
		for _, stat := range pokemon.Stats {
			switch stat.Stat.Name {
			case "hp":
				fmt.Printf(" -hp: %v", stat.BaseStat)
			
			case "attack":
				fmt.Printf("\n -attack: %v", stat.BaseStat)
			
			case "defense":
				fmt.Printf("\n -defense: %v", stat.BaseStat)
			
			case "special-attack":
				fmt.Printf("\n -special-attack: %v", stat.BaseStat)
			
			case "special-defense":
				fmt.Printf("\n -special-defense: %v", stat.BaseStat)
			
			case "speed":
				fmt.Printf("\n -speed: %v\n", stat.BaseStat)
			}
		}

	    fmt.Println("Types:")

		for _, natureType := range pokemon.Types{
			fmt.Printf(" - %v\n", natureType.Type.Name)
		}
			return nil
	}

	fmt.Println("you have not caught that pokemon")

	return nil

}

func commandPokedex(p *pokeapi.Config, name string) error {
	fmt.Println("Your Pokedex:")
	for k, _ := range pokedex{
		fmt.Printf(" - %v\n", k)
	}
	return nil
}

func getCommands() map[string]command{
	
return map[string]command {
	"help" : {
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	},
	"exit" : {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	},
	"map" : {
		name: "map",
		description: "Displays the names of the next 20 location areas in the Pokemon world",
		callback: commandMap,
	},
	"mapb" : {
		name: "mapb",
		description: "Displays the names of the last 20 location areas in the Pokemon world",
		callback: commandMapB,
	},
	"explore" : {
		name: "explore",
		description: "Displays the pokemon that live in the provided location area in the Pokemon world",
		callback: commandExplore,
	},
	"catch" : {
		name: "catch",
		description: "Try to catch a specified pokemon",
		callback: commandCatch,
	},
	"inspect" : {
		name: "inspect",
		description: "Inspect a pokemon you have caught",
		callback: commandInspect,
	},
	"pokedex" : {
		name: "pokedex",
		description: "Inspect a pokemon you have caught",
		callback: commandPokedex,
	},
	
}

}

