package com.tennis.entity;

import java.io.Serializable;
import java.time.LocalDateTime;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

@Entity
@Table(name = "rally")
public class Rally implements Serializable {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "match_id", nullable = false)
    @com.fasterxml.jackson.annotation.JsonProperty("match_id")
    private Long matchId;

    @Column(name = "rally_number", nullable = false)
    @com.fasterxml.jackson.annotation.JsonProperty("rally_number")
    private Integer rallyNumber;

    @Column(nullable = false, length = 20)
    private String scorer;

    @Column(nullable = false, length = 20)
    private String server;

    @Column(nullable = false)
    private LocalDateTime timestamp;

    @Column(name = "video_file", nullable = false, length = 100)
    @com.fasterxml.jackson.annotation.JsonProperty("video_file")
    private String videoFile;

    // Added fields for score snapshot to simplify history/export
    @Column(name = "score_p1")
    @com.fasterxml.jackson.annotation.JsonProperty("score_p1")
    private Integer scoreP1;
    
    @Column(name = "score_p2")
    @com.fasterxml.jackson.annotation.JsonProperty("score_p2")
    private Integer scoreP2;

    public Rally() {
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
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
}
