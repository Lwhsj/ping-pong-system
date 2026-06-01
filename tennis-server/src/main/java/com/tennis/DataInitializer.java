package com.tennis;

import com.tennis.entity.Player;
import com.tennis.repository.PlayerRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

@Component
public class DataInitializer implements CommandLineRunner {
    @Autowired
    private PlayerRepository playerRepository;

    @Override
    public void run(String... args) throws Exception {
        if (playerRepository.count() == 0) {
            playerRepository.save(new Player("Alice", "F", 25));
            playerRepository.save(new Player("Bob", "M", 28));
            playerRepository.save(new Player("Charlie", "M", 22));
            playerRepository.save(new Player("Diana", "F", 24));
            System.out.println("Initialized default players.");
        }
    }
}
