package com.tennis.repository;

import com.tennis.entity.Rally;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.List;

@Repository
public interface RallyRepository extends JpaRepository<Rally, Long> {
    List<Rally> findByMatchId(Long matchId);
}
