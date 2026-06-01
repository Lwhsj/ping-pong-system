package com.tennis.controller;

import com.tennis.dto.RallyDTO;
import com.tennis.entity.Rally;
import com.tennis.service.RallyService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/rally")
@CrossOrigin
public class RallyController {
    @Autowired
    private RallyService rallyService;

    @PostMapping
    public Rally saveRally(@RequestBody RallyDTO rallyDTO) {
        return rallyService.saveRally(rallyDTO);
    }
}
