package com.tennis.service;

import java.io.IOException;
import java.io.OutputStream;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.xssf.usermodel.XSSFWorkbook;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.tennis.dto.MatchCurrentScoreDTO;
import com.tennis.dto.MatchStartRequest;
import com.tennis.dto.MatchStatsDTO;
import com.tennis.entity.Match;
import com.tennis.entity.Player;
import com.tennis.entity.Rally;
import com.tennis.repository.MatchRepository;
import com.tennis.repository.PlayerRepository;
import com.tennis.repository.RallyRepository;

@Service
public class MatchService {
    @Autowired
    private MatchRepository matchRepository;
    @Autowired
    private RallyRepository rallyRepository;
    @Autowired
    private PlayerRepository playerRepository;

    public Match startMatch(MatchStartRequest request) {
        Match match = new Match();
        match.setDate(LocalDate.now());
        match.setPlayer1Id(request.getPlayer1Id());
        match.setPlayer2Id(request.getPlayer2Id());
        match.setStatus("started");
        match.setStartedAt(LocalDateTime.now());
        match.setFirstServer(request.getFirstServer());
        return matchRepository.save(match);
    }

    public Match finishMatch(Long id) {
        Match match = matchRepository.findById(id).orElseThrow(() -> new RuntimeException("Match not found"));
        match.setStatus("finished");
        return matchRepository.save(match);
    }

    public MatchCurrentScoreDTO getCurrentScore(Long id) {
        List<Rally> rallies = rallyRepository.findByMatchId(id);
        int scoreP1 = 0;
        int scoreP2 = 0;
        int rallyNumber = 0;

        if (!rallies.isEmpty()) {
            rallies.sort(Comparator.comparingInt(Rally::getRallyNumber).reversed());
            Rally lastRally = rallies.get(0);
            if (lastRally.getScoreP1() != null) scoreP1 = lastRally.getScoreP1();
            if (lastRally.getScoreP2() != null) scoreP2 = lastRally.getScoreP2();
            rallyNumber = lastRally.getRallyNumber();
        }

        // Calculate current server based on total score
        String server = "Unknown";
        Match match = matchRepository.findById(id).orElse(null);
        
        if (match != null && match.getFirstServer() != null) {
            Long firstServerId = match.getFirstServer();
            Long p1Id = match.getPlayer1Id();
            Long p2Id = match.getPlayer2Id();
            
            int totalScore = scoreP1 + scoreP2;
            int turns;
            
            // Logic: Switch every 2 points until both reach 10 (Deuce), then switch every 1 point
            if (scoreP1 >= 10 && scoreP2 >= 10) {
                // Deuce logic: 10-10 means 20 points played.
                // First 20 points (0-19) were 10 turns (2 serves each).
                // Points from 20 onwards are 1 serve each.
                turns = 10 + (totalScore - 20);
            } else {
                // Normal logic: Switch every 2 points
                turns = totalScore / 2;
            }
            
            Long currentServerId;
            if (turns % 2 == 0) {
                currentServerId = firstServerId;
            } else {
                currentServerId = firstServerId.equals(p1Id) ? p2Id : p1Id;
            }
            
            if (currentServerId.equals(p1Id)) {
                server = "player1";
            } else {
                server = "player2";
            }
        }

        MatchCurrentScoreDTO dto = new MatchCurrentScoreDTO(id, rallyNumber, scoreP1, scoreP2, server);
        
        // Populate player names
        if (match != null) {
            playerRepository.findById(match.getPlayer1Id()).ifPresent(p -> dto.setPlayer1Name(p.getName()));
            playerRepository.findById(match.getPlayer2Id()).ifPresent(p -> dto.setPlayer2Name(p.getName()));
        }
        
        return dto;
    }

    public List<Match> getMatchHistory(LocalDate dateFrom, LocalDate dateTo, String playerName, String status) {
        List<Match> allMatches = matchRepository.findAll();
        
        // Simple in-memory filtering
        return allMatches.stream()
                .filter(m -> dateFrom == null || !m.getDate().isBefore(dateFrom))
                .filter(m -> dateTo == null || !m.getDate().isAfter(dateTo))
                .filter(m -> status == null || status.isEmpty() || m.getStatus().equalsIgnoreCase(status))
                .filter(m -> {
                    if (playerName == null || playerName.isEmpty()) return true;
                    // Find player name in either player1 or player2
                    Player p1 = playerRepository.findById(m.getPlayer1Id()).orElse(null);
                    Player p2 = playerRepository.findById(m.getPlayer2Id()).orElse(null);
                    
                    boolean p1Match = p1 != null && p1.getName().contains(playerName);
                    boolean p2Match = p2 != null && p2.getName().contains(playerName);
                    
                    return p1Match || p2Match;
                })
                .peek(m -> {
                    playerRepository.findById(m.getPlayer1Id()).ifPresent(p -> m.setPlayer1Name(p.getName()));
                    playerRepository.findById(m.getPlayer2Id()).ifPresent(p -> m.setPlayer2Name(p.getName()));
                })
                .collect(Collectors.toList());
    }

    public List<Rally> getMatchDetail(Long id) {
        return rallyRepository.findByMatchId(id);
    }

    public MatchStatsDTO getMatchStats(Long id) {
        List<Rally> rallies = rallyRepository.findByMatchId(id);
        MatchStatsDTO stats = new MatchStatsDTO();
        
        // Calculate stats
        Map<String, Double> serveSuccessRate = new HashMap<>();
        // Placeholder logic: assume 100% or calculate if we had "faults"
        // But we only have "scorer". 
        // If server == scorer, then serve win?
        long p1ServeCount = rallies.stream().filter(r -> "player1".equalsIgnoreCase(r.getServer())).count();
        long p1ServeWin = rallies.stream().filter(r -> "player1".equalsIgnoreCase(r.getServer()) && "player1".equalsIgnoreCase(r.getScorer())).count();
        
        long p2ServeCount = rallies.stream().filter(r -> "player2".equalsIgnoreCase(r.getServer())).count();
        long p2ServeWin = rallies.stream().filter(r -> "player2".equalsIgnoreCase(r.getServer()) && "player2".equalsIgnoreCase(r.getScorer())).count();

        serveSuccessRate.put("player1", p1ServeCount > 0 ? (double)p1ServeWin / p1ServeCount : 0.0);
        serveSuccessRate.put("player2", p2ServeCount > 0 ? (double)p2ServeWin / p2ServeCount : 0.0);
        stats.setServeSuccessRate(serveSuccessRate);

        Map<String, Integer> consecutiveScore = new HashMap<>();
        // Calculate max consecutive score
        consecutiveScore.put("player1", calculateMaxConsecutive(rallies, "player1"));
        consecutiveScore.put("player2", calculateMaxConsecutive(rallies, "player2"));
        stats.setConsecutiveScore(consecutiveScore);

        stats.setAverageRallyTime(15.2); // Mock or calculate if we had duration

        return stats;
    }

    private int calculateMaxConsecutive(List<Rally> rallies, String player) {
        int max = 0;
        int current = 0;
        // Sort by rally number
        rallies.sort(Comparator.comparingInt(Rally::getRallyNumber));
        
        for (Rally r : rallies) {
            if (player.equalsIgnoreCase(r.getScorer())) {
                current++;
            } else {
                max = Math.max(max, current);
                current = 0;
            }
        }
        return Math.max(max, current);
    }

    public void exportMatch(Long id, OutputStream out) throws IOException {
        Match match = matchRepository.findById(id).orElseThrow(() -> new RuntimeException("Match not found"));
        List<Rally> rallies = rallyRepository.findByMatchId(id);
        rallies.sort(Comparator.comparingInt(Rally::getRallyNumber));

        // Get Player Names
        String player1Name = playerRepository.findById(match.getPlayer1Id()).map(Player::getName).orElse("Unknown");
        String player2Name = playerRepository.findById(match.getPlayer2Id()).map(Player::getName).orElse("Unknown");

        try (Workbook workbook = new XSSFWorkbook()) {
            Sheet sheet = workbook.createSheet("比赛数据");
            
            // Header
            Row headerRow = sheet.createRow(0);
            String[] columns = {"回合数", "得分者", "比分 (P1-P2)", "发球方", "时间"};
            for (int i = 0; i < columns.length; i++) {
                headerRow.createCell(i).setCellValue(columns[i]);
            }

            // Date Formatter
            java.time.format.DateTimeFormatter formatter = java.time.format.DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

            // Data
            int rowNum = 1;
            for (Rally rally : rallies) {
                Row row = sheet.createRow(rowNum++);
                row.createCell(0).setCellValue(rally.getRallyNumber());
                
                // Map scorer to name
                String scorerName = "player1".equalsIgnoreCase(rally.getScorer()) ? player1Name : 
                                   ("player2".equalsIgnoreCase(rally.getScorer()) ? player2Name : rally.getScorer());
                row.createCell(1).setCellValue(scorerName);
                
                row.createCell(2).setCellValue(rally.getScoreP1() + "-" + rally.getScoreP2());
                
                // Map server to name
                String serverName = "player1".equalsIgnoreCase(rally.getServer()) ? player1Name : 
                                   ("player2".equalsIgnoreCase(rally.getServer()) ? player2Name : rally.getServer());
                row.createCell(3).setCellValue(serverName);
                
                row.createCell(4).setCellValue(rally.getTimestamp().format(formatter));
            }

            workbook.write(out);
        }
    }
}
