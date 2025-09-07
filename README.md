# Pokedex CLI

A command-line interface (CLI) application for exploring the Pokemon world, catching Pokemon, and managing your collection. Built with Go and powered by the [PokeAPI](https://pokeapi.co/).

## Features

- **Explore Pokemon Locations**: Browse through different location areas in the Pokemon world
- **Catch Pokemon**: Attempt to catch Pokemon with a chance-based system
- **Inspect Pokemon**: View detailed stats and information about caught Pokemon
- **Manage Collection**: View all Pokemon in your Pokedex
- **Caching**: Built-in caching system to improve performance and reduce API calls
- **Interactive CLI**: User-friendly command-line interface

## Commands

| Command   | Description                     | Usage                |
| --------- | ------------------------------- | -------------------- |
| `help`    | Display help message            | `help`               |
| `map`     | Show next 20 location areas     | `map`                |
| `mapb`    | Show previous 20 location areas | `mapb`               |
| `explore` | Explore Pokemon in a location   | `explore <location>` |
| `catch`   | Attempt to catch a Pokemon      | `catch <pokemon>`    |
| `inspect` | Inspect a caught Pokemon        | `inspect <pokemon>`  |
| `pokedex` | List all caught Pokemon         | `pokedex`            |
| `exit`    | Exit the application            | `exit`               |

## Installation

### Prerequisites

- Go 1.25.0 or later
- Internet connection (for API calls)

### Build and Run

1. Clone the repository:

```bash
git clone https://github.com/abdo-355/pokedex.git
cd pokedex
```

2. Build the application:

```bash
go build
```

3. Run the application:

```bash
./pokedex
```

Or run directly with Go:

```bash
go run main.go
```

## Usage Example

```
Pokedex > help
Welcome to the Pokedex!
Usage:

help: Displays a help message
map: displays the names of 20 location areas in the Pokemon world.
mapb: displays the names of the previous 20 location areas in the Pokemon world.
explore <location>: explore the pokemon available at a certain location
catch <pokemon>: catch a pokemon and add it to your collection
inspect <pokemon>: inspect a pokemon in your collection
exit: Exit the Pokedex

Pokedex > map
pallet-town
viridian-city
pewter-city
...

Pokedex > explore pallet-town
Exploring pallet-town...
Found Pokemon:
- rattata
- pidgey

Pokedex > catch rattata
Throwing a Pokeball at rattata...
rattata was caught!
You may now inspect it with the inspect command.

Pokedex > inspect rattata
Name: rattata
Height: 3
Weight: 35
Stats:
  -hp: 30
  -attack: 56
  -defense: 35
  -special-attack: 25
  -special-defense: 35
  -speed: 72
Types:
  - normal

Pokedex > pokedex
Your Pokedex:
 - rattata
```

## Project Structure

```
pokedex/
├── main.go                 # Main application entry point
├── go.mod                  # Go module file
├── repl_test.go           # Tests for REPL functionality
└── internal/
    └── pokecache/
        ├── main.go        # Cache implementation
        └── main_test.go   # Cache tests
```

## Architecture

### Core Components

1. **CLI Interface**: Interactive command-line interface with command parsing
2. **API Client**: HTTP client for interacting with PokeAPI
3. **Cache System**: Thread-safe caching mechanism with automatic expiration
4. **Pokemon Management**: Collection system for caught Pokemon

### Cache System

The application includes a sophisticated caching system (`internal/pokecache/`) that:

- Stores API responses to reduce network calls
- Automatically expires entries after a configurable duration (7 seconds by default)
- Uses goroutines for background cleanup
- Provides thread-safe concurrent access

### Pokemon Catching Mechanics

The catching system uses a probability-based approach:

- Higher base experience Pokemon are harder to catch
- Success rate is calculated based on the Pokemon's base experience
- Already caught Pokemon are automatically successful

## Testing

Run the test suite:

```bash
go test ./...
```

The project includes comprehensive tests for:

- Cache functionality (add, get, expiration, concurrent access)
- Input parsing and cleaning
- Command execution

## API Integration

This application integrates with the [PokeAPI](https://pokeapi.co/), a free RESTful API that provides comprehensive Pokemon data. The app makes requests to:

- Location areas: `https://pokeapi.co/api/v2/location-area/`
- Pokemon data: `https://pokeapi.co/api/v2/pokemon/`

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- [PokeAPI](https://pokeapi.co/) for providing the Pokemon data
- The Pokemon Company for creating the amazing Pokemon universe
