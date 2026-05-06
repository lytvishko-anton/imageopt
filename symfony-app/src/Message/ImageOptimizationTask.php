<?php
namespace App\Message;

class ImageOptimizationTask {
    public function __construct(
        private string $filePath
    ) {}

    public function getFilePath(): string {
        return $this->filePath;
    }
}