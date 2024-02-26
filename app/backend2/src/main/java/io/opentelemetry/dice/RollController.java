package io.opentelemetry.dice;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.Optional;


@RestController
public class RollController {

	private static final Logger logger = LogManager.getLogger(RollController.class);


	@GetMapping("/rolldice")
	public String index(@RequestParam("player") Optional<String> player) {
		int result = this.getRandomNumber(-2, 6);
		if(result < 1) {
			logger.warn("Illegal number rolled, setting result to '1'");
			result = 1;
		}
		if(player.isPresent()) {
			logger.info(player.get() + " is rolling the dice: " + result);
		} else {
			logger.info("Anonymous player is rolling the dice: " + result);
		}
		return Integer.toString(result);
	}

    public int getRandomNumber(int min, int max) {
        return (int) ((Math.random() * (max - min)) + min);
    }
}