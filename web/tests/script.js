print("Hello, world!");
setInterval(() => {
    let contestants = getContestantsByRanking("survived");
    for (let i = 0; i < contestants.length; i++) {
        for (let j = i + 1; j < contestants.length; j++) {
            print("create!");
            createMatch([contestants[i], contestants[j]]);
        }
    }
}, 5 * 1000);

function onAiAssigned(contestant) {

}


function onMatchFinished(players, tag, replay) {
    players.forEach(player => {
        let newPoints = player.contestant.points + player.score;
        updateContestant(player.contestant, { "points": newPoints });
    });
}

//转义后
//print(\"Hello, world!\");\nsetInterval(() => {\n    let contestants = getContestantsByRanking(\"survived\");\n    for (let i = 0; i < contestants.length; i++) {\n        for (let j = i + 1; j < contestants.length; j++) {\n            print(\"create!\");\n            createMatch([contestants[i], contestants[j]]);\n        }\n    }\n}, 5 * 1000);\n\nfunction onAiAssigned(contestant) {\n\n}\n\n\nfunction onMatchFinished(players, tag, replay) {\n    players.forEach(player => {\n        let newPoints = player.contestant.points + player.score;\n        updateContestant(player.contestant, { \"points\": newPoints });\n    });\n}