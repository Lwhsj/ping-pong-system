package com.tennis.controller;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.util.List;

import org.springframework.core.io.Resource;
import org.springframework.core.io.UrlResource;
import org.springframework.core.io.support.ResourceRegion;
import org.springframework.http.HttpRange;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/api")
@CrossOrigin
public class VideoController {

    private final Path fileStorageLocation;

    public VideoController() {
        this.fileStorageLocation = Paths.get("uploads").toAbsolutePath().normalize();
        try {
            Files.createDirectories(this.fileStorageLocation);
        } catch (Exception ex) {
            throw new RuntimeException("Could not create the directory where the uploaded files will be stored.", ex);
        }
    }

    @PostMapping("/upload/video")
    public ResponseEntity<String> uploadVideo(@RequestParam("file") MultipartFile file,
                                              @RequestParam("matchId") Long matchId) {
        try {
            if (file.isEmpty()) {
                return ResponseEntity.badRequest().body("File is empty");
            }

            // Normalize file name
            String fileName = file.getOriginalFilename();
            if (fileName == null) {
                fileName = "video_" + matchId + "_" + System.currentTimeMillis() + ".webm";
            }
            
            // Clean filename to avoid security issues (simple check)
            fileName = fileName.replaceAll("[^a-zA-Z0-9._-]", "_");

            Path targetLocation = this.fileStorageLocation.resolve(fileName);
            Files.copy(file.getInputStream(), targetLocation, StandardCopyOption.REPLACE_EXISTING);

            return ResponseEntity.ok(fileName);
        } catch (IOException ex) {
            return ResponseEntity.internalServerError().body("Could not upload file: " + ex.getMessage());
        }
    }

    @GetMapping("/video/{fileName:.+}")
    public ResponseEntity<ResourceRegion> getVideo(@PathVariable String fileName, @RequestHeader(value = "Range", required = false) String rangeHeader) {
        try {
            Path filePath = this.fileStorageLocation.resolve(fileName).normalize();
            Resource resource = new UrlResource(filePath.toUri());

            if (resource.exists()) {
                long contentLength = resource.contentLength();
                ResourceRegion region;
                
                // Determine Content-Type
                String contentType = null;
                try {
                    contentType = Files.probeContentType(filePath);
                } catch (IOException ex) {
                    // Ignored
                }
                
                // Fallback for common video types if probe fails
                if (contentType == null) {
                    String lowerFileName = fileName.toLowerCase();
                    if (lowerFileName.endsWith(".webm")) {
                        contentType = "video/webm";
                    } else if (lowerFileName.endsWith(".mp4")) {
                        contentType = "video/mp4";
                    } else {
                        contentType = "application/octet-stream";
                    }
                }
                
                MediaType mediaType = MediaType.parseMediaType(contentType);

                if (rangeHeader != null) {
                    long start = 0;
                    long end = contentLength - 1;
                    
                    // Parse Range Header manually or using HttpRange
                    // Example rangeHeader: "bytes=0-" or "bytes=0-1000"
                    List<HttpRange> ranges = HttpRange.parseRanges(rangeHeader);
                    if (!ranges.isEmpty()) {
                        HttpRange range = ranges.get(0);
                        start = range.getRangeStart(contentLength);
                        end = range.getRangeEnd(contentLength);
                    }
                    
                    // Calculate length of the chunk
                    long rangeLength = Math.min(1024 * 1024, end - start + 1); // 1MB chunk size
                    region = new ResourceRegion(resource, start, rangeLength);
                    
                    return ResponseEntity.status(HttpStatus.PARTIAL_CONTENT)
                            .contentType(mediaType)
                            .body(region);
                } else {
                    // No range, return full content (still using ResourceRegion for consistency)
                    region = new ResourceRegion(resource, 0, contentLength);
                    return ResponseEntity.status(HttpStatus.OK)
                            .contentType(mediaType)
                            .body(region);
                }
            } else {
                return ResponseEntity.notFound().build();
            }
        } catch (Exception ex) {
            return ResponseEntity.badRequest().build();
        }
    }
}
