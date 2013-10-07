# design notes

Canvas mode will initially be the only mode. I'm picturing a pretty
simple canvas model where nodes are placed on a strict grid, and
connections run in fixed lanes between the grids. I'm not going to
care much about collisions etc, and connections will be color-coded to
show which nodes connect to what. One nice thing would be to highlight
the complete chain from source to sink when selecting a source.

# Model

OK, so here are the node types:

* Source: Track

  A track is a source. A track is a stream of sound with an indefinite
  extent. Into a track can be loaded a WAV file, for
  example. Eventually, full editing of the contents of a track (with
  fades etc) plus track envelopes like panning and volume should be
  added, but for the first prototype, a track is just a time-shifted
  sample.

* FX: Playback delay

  Simply insert silence into the stream until a set time, when it
  starts playing back its input.

* FX: Noise gate

  Mutes sound below a certain dB. Level needs to be configurable.

* FX: Compressor

  Well, I'll need to read up on how FFT works in general, but a noise
  gate and a compressor should be implementable as the same code path
  with different inputs, if I remember correctly.

* FX: EQ

  Probably just a very simple EQ model to start with.

* FX: Splitter

  Single input, multiple outputs.

* FX: Mixer

  Multiple inputs, single output

* Sink: Speaker

  Plays its incoming stream to the speakers / headphones.

* Sink: WAV

  Writes its incoming stream to a WAV file on disk.

# UI

* Canvas

* Box

* Link

* PopupMenu

* Dialog

* Slider

* NumberView

* Track

* Block/Clip - in a track


# Track view

Vertically, tracks fill the track view, one track for each
source. Sources can be positioned in time relative to each other. All
this does is time-shift the source.

Future: Full DAW-style track editing of sources, with looping,
cutting, mixing etc. Each track is a source that can be linked in the
canvas view.

# Input controls

* F1: Canvas view
* F2: Track view

