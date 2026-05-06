<?php

namespace App\Controller;

use App\Message\ImageOptimizationTask;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\BinaryFileResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Messenger\MessageBusInterface;
use Symfony\Component\Routing\Annotation\Route;

class ImageUploadController extends AbstractController
{
    #[Route('/upload', name: 'app_image_upload', methods: ['POST'])]
    public function upload(Request $request, MessageBusInterface $bus): Response
    {
        $file = $request->files->get('image');
        
        if (!$file) {
            return $this->json(['error' => 'No file received'], 400);
        }

        $uploadDir = '/shared_uploads';
        $originalName = $file->getClientOriginalName();
        $targetPath = $uploadDir . '/' . $originalName;

        // Use pathinfo to get the filename without extension for the WebP name
        $baseName = pathinfo($originalName, PATHINFO_FILENAME);
        $webpName = $baseName . '.webp';

        if (copy($file->getPathname(), $targetPath)) {
            // Dispatch the task to the Go worker via RabbitMQ
            $bus->dispatch(new ImageOptimizationTask($targetPath));
            
            return $this->json([
                'status' => 'Processing',
                'webpName' => $webpName 
            ]);
        }

        return $this->json(['error' => 'Failed to save original file'], 500);
    }

    #[Route('/check-status/{filename}', name: 'app_check_status', methods: ['GET'])]
    public function checkStatus(string $filename): Response
    {
        $path = '/shared_uploads/' . $filename;
        if (!file_exists($path)) {
            return new Response("No", 404);
        }
        return new Response("Yes", 200);
    }

    #[Route('/downloads/{filename}', name: 'app_download', methods: ['GET'])]
    public function download(string $filename): Response
    {
        $path = '/shared_uploads/' . $filename;
        if (!file_exists($path)) return new Response("File gone", 404);

        $response = new \Symfony\Component\HttpFoundation\BinaryFileResponse($path);
        $response->deleteFileAfterSend(true);
        return $response;
    }
}