package main

var IMPORTS = map[string]string{

    "scale": `
        :ut 220;
        :A  [ut]Hz;
        :Bb [ut#]Hz;
        :B  [ut##]Hz;
        :C  [ut###]Hz;
        :Db [ut####]Hz;
        :D  [ut#####]Hz;
        :Eb [ut######]Hz; 
        :E  [ut#######]Hz;
        :F  [ut########]Hz;
        :Gb [ut#########]Hz;
        :G  [ut##########]Hz;
        :Ab [ut###########]Hz;
    `,

    "timing": `
        :beats (loop-count freq -- beat-age beat-num)
            1 dmod rot mod
        ;
    `,
}
