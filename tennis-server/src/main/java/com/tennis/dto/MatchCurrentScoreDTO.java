package com.tennis.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

public class MatchCurrentScoreDTO {
    @JsonProperty("match_id")
    private Long matchId;

    @JsonProperty("rally_number")
    private Integer rallyNumber;

    @JsonProperty("score_p1")
    private Integer scoreP1;

    @JsonProperty("score_p2")
    private Integer scoreP2;

    private String server;

    @JsonProperty("player1_name")
    private String player1Name;

    @JsonProperty("player2_name")
    private String player2Name;

    public MatchCurrentScoreDTO() {
    }

    public MatchCurrentScoreDTO(Long matchId, Integer rallyNumber, Integer scoreP1, Integer scoreP2, String server) {
        this.matchId = matchId;
        this.rallyNumber = rallyNumber;
        this.scoreP1 = scoreP1;
        this.scoreP2 = scoreP2;
        this.server = server;
    }

    public Long getMatchId() {
        return matchId;
    }

    public void setMatchId(Long matchId) {
        this.matchId = matchId;
    }

    public Integer getRallyNumber() {
        return rallyNumber;
    }

    public void setRallyNumber(Integer rallyNumber) {
        this.rallyNumber = rallyNumber;
    }

    public Integer getScoreP1() {
        return scoreP1;
    }

    public void setScoreP1(Integer scoreP1) {
        this.scoreP1 = scoreP1;
    }

    public Integer getScoreP2() {
        return scoreP2;
    }

    public void setScoreP2(Integer scoreP2) {
        this.scoreP2 = scoreP2;
    }

    public String getServer() {
        return server;
    }

    public void setServer(String server) {
        this.server = server;
    }

    public String getPlayer1Name() {
        return player1Name;
    }

    public void setPlayer1Name(String player1Name) {
        this.player1Name = player1Name;
    }

    public String getPlayer2Name() {
        return player2Name;
    }

    public void setPlayer2Name(String player2Name) {
        this.player2Name = player2Name;
    }
}
