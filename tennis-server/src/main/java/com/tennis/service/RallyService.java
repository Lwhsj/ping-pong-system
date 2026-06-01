package com.tennis.service;

import com.tennis.dto.RallyDTO;
import com.tennis.entity.Rally;
import com.tennis.repository.RallyRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import java.util.List;
import java.util.Comparator;

@Service
public class RallyService {
    @Autowired
    private RallyRepository rallyRepository;

    public Rally saveRally(RallyDTO rallyDTO) {
        Rally rally = new Rally();
        rally.setMatchId(rallyDTO.getMatchId());
        rally.setRallyNumber(rallyDTO.getRallyNumber());
        rally.setScorer(rallyDTO.getScorer());
        rally.setServer(rallyDTO.getServer());
        rally.setTimestamp(rallyDTO.getTimestamp());
        rally.setVideoFile(rallyDTO.getVideoFile());

        // Calculate score
        // Get previous score
        List<Rally> rallies = rallyRepository.findByMatchId(rally.getMatchId());
        int scoreP1 = 0;
        int scoreP2 = 0;
        
        if (!rallies.isEmpty()) {
            // Sort by rally number descending
            rallies.sort(Comparator.comparingInt(Rally::getRallyNumber).reversed());
            Rally lastRally = rallies.get(0);
            if (lastRally.getScoreP1() != null) scoreP1 = lastRally.getScoreP1();
            if (lastRally.getScoreP2() != null) scoreP2 = lastRally.getScoreP2();
        }

        if ("player1".equalsIgnoreCase(rallyDTO.getScorer())) {
            scoreP1++;
        } else if ("player2".equalsIgnoreCase(rallyDTO.getScorer())) {
            scoreP2++;
        }

        rally.setScoreP1(scoreP1);
        rally.setScoreP2(scoreP2);

        return rallyRepository.save(rally);
    }

    public List<Rally> getRalliesByMatchId(Long matchId) {
        return rallyRepository.findByMatchId(matchId);
    }
}
