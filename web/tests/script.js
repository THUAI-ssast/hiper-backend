print("Hello, world!");
setInterval(() => {
    print("aaa")
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