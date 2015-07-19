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

    "divisions": `
        :divisions (loop-count freq -- beat-age beat-num)

            1 dmod (-- loop-count beat-age beat-count)
            rot    (-- beat-age beat-count loop-count)
            mod

        ;
    `,

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
