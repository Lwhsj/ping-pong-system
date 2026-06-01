package com.tennis.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.Map;

public class MatchStatsDTO {
    @JsonProperty("serve_success_rate")
    private Map<String, Double> serveSuccessRate;

    @JsonProperty("consecutive_score")
    private Map<String, Integer> consecutiveScore;

    @JsonProperty("average_rally_time")
    private Double averageRallyTime;

    public Map<String, Double> getServeSuccessRate() {
        return serveSuccessRate;
    }

    public void setServeSuccessRate(Map<String, Double> serveSuccessRate) {
        this.serveSuccessRate = serveSuccessRate;
    }

    public Map<String, Integer> getConsecutiveScore() {
        return consecutiveScore;
    }

    public void setConsecutiveScore(Map<String, Integer> consecutiveScore) {
        this.consecutiveScore = consecutiveScore;
    }

    public Double getAverageRallyTime() {
        return averageRallyTime;
    }

    public void setAverageRallyTime(Double averageRallyTime) {
        this.averageRallyTime = averageRallyTime;
    }
}
