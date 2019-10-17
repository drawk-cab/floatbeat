package main

var IMPORTS = map[string]string{

    "scale": `
        :ut 220;
        :A  [ut]Hz;
        :Bb [ut ♯]Hz;
        :B  [ut ♯♯]Hz;
        :C  [ut ♯♯♯]Hz;
        :Db [ut ♯♯♯♯]Hz;
        :D  [ut ♯♯♯♯♯]Hz;
        :Eb [ut ♯♯♯♯♯♯]Hz; 
        :E  [ut ↑ ♭♭♭♭♭]Hz;
        :F  [ut ↑ ♭♭♭♭]Hz;
        :Gb [ut ↑ ♭♭♭]Hz;
        :G  [ut ↑ ♭♭]Hz;
        :Ab [ut ↑ ♭]Hz;
    `,

    "c_just_scale": `
        :ut [220 5/3*] ; (tonic is C)

        :A     220 Hz ;             :Dha A;
        :C     [ut] Hz ;            :Sa C;
        :Db    [ut 24/25*] Hz ;
        :D     [ut 8/9*] Hz   ;     :Re D;
        :Eb    [ut 5/6*] Hz   ;
        :E     [ut 4/5*] Hz   ;     :Ga E;
        :F     [ut 3/4*] Hz   ;     :Ma F;
        :Gb    [ut 32/45*] Hz ;
        :G     [ut 2/3*] Hz   ;     :Pa G;
        :Ab    [ut 5/8*] Hz   ;
        :Bb    [ut 5/9*] Hz   ;
        :B     [ut 8/15*] Hz  ;     :Ni B;
    `,

    "g_just_scale": `
        :ut [220 5/9*] ; (tonic is G)

        :A     440 Hz ;
        :G     [ut] Hz ;
        :Ab    [ut 24/25*] Hz ;
        :Bb    [ut 5/6*] Hz   ;
        :B     [ut 4/5*] Hz   ;
        :C     [ut 3/4*] Hz   ;
        :Db    [ut 32/45*] Hz ;
        :D     [ut 2/3*] Hz   ;
        :Eb    [ut 5/8*] Hz   ;
        :E     [ut 3/5*] Hz   ;
        :F     [ut 5/9*] Hz   ;
        :Gb    [ut 8/15*] Hz  ;
    `,

    "ticks": `
:ticks (loop-length freq -- tick-in-loop tick-age)
  dup T *    (-- loop-length freq time-in-ticks)
  rot        (-- freq time-in-ticks loop-length)
  swap       (-- freq loop-length time-in-ticks)
  1 dmod     (-- freq loop-length tick-fraction tick-count)
  rot        (-- freq tick-fraction tick-count loop-length)
  mod        (-- freq tick-fraction tick-in-loop)
  hide       (-- tick-in-loop freq tick-fraction)
  \          (-- tick-in-loop tick-age)
;`,

    "env": `
:env (signal note-age attack-time decay-time -- new-signal)
  rot ddup                            (-- signal attack-time decay-time note-age decay-time note-age)
  < if
     (note-age) drop
     (decay-time) drop
     (attack-time) drop
     (signal) drop
     0
  else
  rot ddup                            (-- signal decay-time note-age attack-time note-age attack-time)
  > if                                (-- signal decay-time note-age attack-time)
     (note-age attack-time) -         (-- signal decay-time segment-age)
     (decay-time segment-age) \ 1~    (-- signal amplitude)
     (signal amplitude) *
  else
     rot                              (-- signal note-age attack-time decay-time)
     (decay-time) drop
     (note-age attack-time) /         (-- signal amplitude)
     (signal amplitude) *
  then then
;`,

    "lowpass": `
        :lowpass (input alpha -- output)
            swap _lowpass delta -
            swap *
            _lowpass delta +
            dup keep _lowpass
        ;
    `,

    "highpass": `
        :highpass (input alpha -- output)
            swap dup keep _highpass_in
            _highpass_in delta -
            _highpass_out delta +
            * dup keep _highpass_out
        ;
    `,
}
