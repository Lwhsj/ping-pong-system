package com.tennis.dto;

import java.time.LocalDateTime;

import com.fasterxml.jackson.annotation.JsonProperty;

public class RallyDTO {
    @JsonProperty("match_id")
    private Long matchId; // Request uses string "1"

    @JsonProperty("set_number")
    private Integer setNumber;

    @JsonProperty("rally_number")
    private Integer rallyNumber;

    private String scorer;
    private String server;
    private LocalDateTime timestamp;

    @JsonProperty("video_file")
    private String videoFile;

    // Getters and setters
    public Long getMatchId() {
        return matchId;
    }

    public void setMatchId(Long matchId) {
        this.matchId = matchId;
    }

    public Integer getSetNumber() {
        return setNumber;
    }

    public void setSetNumber(Integer setNumber) {
        this.setNumber = setNumber;
    }

    public Integer getRallyNumber() {
        return rallyNumber;
    }

    public void setRallyNumber(Integer rallyNumber) {
        this.rallyNumber = rallyNumber;
    }

    public String getScorer() {
        return scorer;
    }

    public void setScorer(String scorer) {
        this.scorer = scorer;
    }

    public String getServer() {
        return server;
    }

    public void setServer(String server) {
        this.server = server;
    }

    public LocalDateTime getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(LocalDateTime timestamp) {
        this.timestamp = timestamp;
    }

    public String getVideoFile() {
        return videoFile;
    }

    public void setVideoFile(String videoFile) {
        this.videoFile = videoFile;
    }
}
