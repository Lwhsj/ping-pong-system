package com.tennis.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

public class MatchStartRequest {
    @JsonProperty("player1_id")
    private Long player1Id;

    @JsonProperty("player2_id")
    private Long player2Id;

    @JsonProperty("first_server")
    private Long firstServer;

    public Long getPlayer1Id() {
        return player1Id;
    }

    public void setPlayer1Id(Long player1Id) {
        this.player1Id = player1Id;
    }

    public Long getPlayer2Id() {
        return player2Id;
    }

    public void setPlayer2Id(Long player2Id) {
        this.player2Id = player2Id;
    }

    public Long getFirstServer() {
        return firstServer;
    }

    public void setFirstServer(Long firstServer) {
        this.firstServer = firstServer;
    }
}
