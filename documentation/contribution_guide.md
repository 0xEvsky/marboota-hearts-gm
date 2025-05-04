# Contribution Guide

## Getting Started

This guide will help you understand how to contribute to the Marboota Game project, a multiplayer card game built with Godot 4.

## Prerequisites

- [Godot 4.3](https://godotengine.org/download) or newer
- Basic understanding of GDScript
- Git for version control
- Access to the project repository

## Project Setup

1. Clone the repository:
   ```bash
   git clone [repository-url]
   ```

2. Open Godot Engine and import the project by selecting the `project.godot` file in the `game` directory.

3. Run the project by pressing F5 or clicking the Play button in the editor.

## Project Structure

```
game/
├── assets/             # Game assets (graphics, UI, fonts)
├── scenes/             # Scene files (.tscn)
├── scripts/            # GDScript files (.gd)
├── .gitattributes      # Git configuration
├── .gitignore          # Files to ignore in version control
├── export_presets.cfg  # Export configurations
├── icon.svg            # Project icon
├── project.godot       # Godot project file
```

## Development Workflow

### Running the Game

To test the game with networking:

1. Start the backend server (not included in the uploaded files)
2. Run the game from Godot editor
3. The game will connect to `ws://localhost:3000/ws` by default

### Making Changes

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes in the Godot editor or by editing script files directly.

3. Test your changes by running the game.

4. Commit your changes with descriptive messages:
   ```bash
   git commit -m "Add feature: description of what you added"
   ```

5. Push your branch to the repository:
   ```bash
   git push origin feature/your-feature-name
   ```

6. Create a pull request to merge your changes into the main branch.

## Coding Guidelines

### GDScript Style

- Use PascalCase for class names: `PlayerManager`
- Use snake_case for variables and functions: `player_count`, `update_position()`
- Add type hints where possible: `var speed: float = 10.0`
- Use meaningful variable and function names
- Include comments for complex logic

### Scene Organization

- Create reusable components as separate scenes
- Use node groups to organize related nodes
- Follow the existing scene hierarchy patterns

## Common Contribution Tasks

### Adding a New Feature

1. Identify where your feature fits in the existing architecture
2. Create new script files if needed
3. Integrate with existing systems (NetworkManager, EventManager, etc.)
4. Add appropriate signals for communication between components
5. Update the UI if your feature requires user interaction
6. Add documentation for your new feature

### Fixing Bugs

1. Reproduce the bug to understand its cause
2. Locate the relevant code in the project
3. Fix the issue with minimal changes
4. Add tests if possible to prevent regression
5. Update documentation if your fix changes functionality

### Improving Existing Systems

1. Understand the current implementation thoroughly
2. Make incremental improvements rather than complete rewrites
3. Ensure backward compatibility with existing code
4. Update documentation to reflect your changes

## Key Systems to Be Aware Of

When contributing, be aware of these core systems:

1. **NetworkManager**: Handles WebSocket communication with the server
2. **EventManager**: Manages game events and request/response handling
3. **PlayerManager**: Controls player creation, positioning, and state
4. **Seat System**: Manages player seating arrangements
5. **Ready System**: Handles player readiness state

## Testing Your Changes

- Test network functionality with both single and multiple clients
- Verify UI elements work correctly at different screen resolutions
- Check that your changes don't break existing functionality
- Test edge cases, especially for network-related features

## Documentation

When adding new features or making significant changes:

1. Update relevant markdown files in the docs directory
2. Add inline comments for complex logic
3. Update function documentation for public methods
4. Consider adding examples for usage of new features

## Getting Help

If you need assistance while contributing:

- Check the existing documentation
- Look at similar implementations in the codebase
- Reach out to project maintainers

Thank you for contributing to the Marboota Game project!