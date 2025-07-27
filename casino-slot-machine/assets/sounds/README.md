# Sound Assets

This directory contains audio files for the slot machine game. The current implementation uses Web Audio API to generate synthetic sounds, but actual audio files can be placed here for enhanced audio experience.

## File Structure

- `spin.ogg` - Reel spinning sound effect
- `win-small.ogg` - Small win celebration sound
- `win-big.ogg` - Big win celebration sound  
- `win-jackpot.ogg` - Jackpot fanfare sound
- `reel-stop.ogg` - Individual reel stopping sound
- `button-click.ogg` - UI button click sound
- `coin-drop.ogg` - Coin collection sound
- `background.ogg` - Optional background music (looping)

## Audio Implementation

The game uses a hybrid approach:
1. **Synthetic Audio**: Web Audio API generates sounds programmatically for guaranteed compatibility
2. **Preloaded Files**: Optional enhancement with actual audio files for richer experience
3. **Fallback System**: Gracefully degrades if audio files are unavailable

## Browser Compatibility

- Web Audio API sounds work in all modern browsers
- Audio files require user interaction to start (autoplay policies)
- Volume controls and mute functionality included
- Settings persist via localStorage
