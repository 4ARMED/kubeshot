package screenshot

import (
	"net/url"
	"strings"
	"sync"

	"github.com/4armed/kubeshot/internal/config"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/kubicorn/kubicorn/pkg/logger"
)

// Process works through the supplied URLs calling toPng on each
func Process(urls []string, c *config.Config) {
	logger.Info("processing %v URLs with %d workers", len(urls), c.NumberOfWorkers)

	// Workers get URLs from this channel
	urlsToProcess := make(chan string)

	var wg sync.WaitGroup
	for w := 0; w < c.NumberOfWorkers; w++ {
		logger.Debug("Setting up worker %d", w)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlsToProcess {
				screenshot(url, c)
			}
		}()
	}

	// Feed workers with URLs
	go func() {
		// Workers will exit from range loop when channel is closed
		defer close(urlsToProcess)
		for _, u := range urls {
			urlsToProcess <- u
		}
	}()

	logger.Info("Waiting for workers...")
	wg.Wait()
}

func screenshot(url string, c *config.Config) {
	logger.Debug("processing URL %s", url)
	filename, err := getFilename(url, c.OutputDir)
	if err != nil {
		logger.Warning("getFilename returned error: %v", err)
	}

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	// Ignore certificate warnings
	browser.MustIgnoreCertErrors(true)

	page, err := browser.Page(proto.TargetCreateTarget{
		URL: url,
	})
	if err != nil {
		logger.Warning("could not load page %s: %v", url, err)
		return
	}

	err = page.WaitLoad()
	if err != nil {
		logger.Warning("could not load URL %s: %v", url, err)
		return
	}

	img, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatPng,
		Quality: 90,
	})
	if err != nil {
		logger.Warning("screenshot returned error: %v", err)
	}
	_ = utils.OutputFile(filename, img)
}

// getFilename returns the full path to the calculated filename
func getFilename(targetURL string, directory string) (string, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		logger.Warning("could not parse URL %v", targetURL)
	}
	filename := strings.ReplaceAll(parsedURL.Host+".png", ":", "-")
	return directory + filename, err
}
