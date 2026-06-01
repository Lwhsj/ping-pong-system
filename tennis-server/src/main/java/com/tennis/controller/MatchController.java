package com.tennis.controller;

import com.tennis.dto.MatchCurrentScoreDTO;
import com.tennis.dto.MatchStartRequest;
import com.tennis.dto.MatchStatsDTO;
import com.tennis.entity.Match;
import com.tennis.entity.Rally;
import com.tennis.service.MatchService;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.time.LocalDate;
import java.util.List;

@RestController
@RequestMapping("/api")
@CrossOrigin
public class MatchController {
    @Autowired
    private MatchService matchService;

    @PostMapping("/match/start")
    public Match startMatch(@RequestBody MatchStartRequest request) {
        return matchService.startMatch(request);
    }

    @PostMapping("/match/{id}/finish")
    public Match finishMatch(@PathVariable Long id) {
        return matchService.finishMatch(id);
    }

    @GetMapping("/match/{id}/current")
    public MatchCurrentScoreDTO getCurrentScore(@PathVariable Long id) {
        return matchService.getCurrentScore(id);
    }

    @GetMapping("/matches")
    public List<Match> getMatchHistory(
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate dateFrom,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate dateTo,
            @RequestParam(required = false) String player,
            @RequestParam(required = false) String status) {
        return matchService.getMatchHistory(dateFrom, dateTo, player, status);
    }

    @GetMapping("/match/{id}/detail")
    public List<Rally> getMatchDetail(@PathVariable Long id) {
        return matchService.getMatchDetail(id);
    }

    @GetMapping("/match/{id}/stats")
    public MatchStatsDTO getMatchStats(@PathVariable Long id) {
        return matchService.getMatchStats(id);
    }

    @GetMapping("/match/{id}/export")
    public void exportMatch(@PathVariable Long id, HttpServletResponse response) throws IOException {
        response.setContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet");
        response.setHeader("Content-Disposition", "attachment; filename=match_" + id + ".xlsx");
        matchService.exportMatch(id, response.getOutputStream());
    }
}
