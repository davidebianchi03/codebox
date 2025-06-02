package randomnames

import (
	"math/rand"
	"strings"
	"time"
)

var adjectives = []string{
	"admiring", "adoring", "agitated", "amazing", "angry", "awesome", "blissful",
	"bold", "boring", "brave", "busy", "charming", "clever", "cool", "compassionate",
	"competent", "condescending", "confident", "cranky", "crazy", "dazzling",
	"determined", "distracted", "dreamy", "eager", "ecstatic", "elated", "elegant",
	"eloquent", "epic", "fervent", "frosty", "gallant", "gifted", "goofy", "gracious",
	"happy", "hardcore", "heuristic", "hopeful", "hungry", "infallible", "inspiring",
	"jolly", "jovial", "keen", "kind", "laughing", "loving", "lucid", "mystifying",
	"modest", "musing", "naughty", "nervous", "nice", "nifty", "nostalgic", "objective",
	"optimistic", "peaceful", "pedantic", "pensive", "practical", "priceless", "quirky",
	"quizzical", "recursing", "relaxed", "reverent", "romantic", "sad", "serene",
	"sharp", "silly", "sleepy", "stoic", "strange", "stupefied", "suspicious", "sweet",
	"tender", "thirsty", "trusting", "unruffled", "upbeat", "vibrant", "vigilant",
	"vigorous", "wizardly", "wonderful", "xenodochial", "youthful", "zealous", "zen",
}

var surnames = []string{
	"agnesi", "archimedes", "ardinghelli", "babbage", "banach", "bardeen", "bartik",
	"bassi", "bell", "benz", "bhabha", "bhaskara", "black", "blackburn", "blackwell",
	"bohr", "booth", "borg", "bose", "boyd", "brahmagupta", "brattain", "brown",
	"carson", "cartwright", "chandrasekhar", "chaplygin", "chatelet", "clarke", "colden",
	"cori", "cray", "curie", "darwin", "davinci", "dijkstra", "dubinsky", "easley",
	"edison", "einstein", "elion", "engelbart", "euclid", "euler", "fermat", "fermi",
	"feynman", "franklin", "galileo", "galois", "ganguly", "gates", "goldberg",
	"goldstine", "goldwasser", "golick", "goodall", "grothendieck", "haibt", "hamilton",
	"haslett", "hawking", "hellman", "heisenberg", "hermann", "herschel", "hertz",
	"heyrovsky", "hodgkin", "hoover", "hopper", "hugle", "hypatia", "ishizaka", "jackson",
	"jang", "jennings", "jepsen", "johnson", "joliot", "jones", "kalam", "keller",
	"khorana", "kilby", "kirch", "knuth", "kowalevski", "lalande", "lamarr", "lamé",
	"leakey", "leavitt", "lederberg", "lehmann", "lewin", "lichterman", "liskov",
	"lovelace", "lumiere", "mahavira", "margulis", "matsumoto", "maxwell", "mayer",
	"mccarthy", "mcclintock", "mclaren", "mclean", "mcnulty", "meitner", "meninsky",
	"mestorf", "mirzakhani", "morse", "newton", "nightingale", "nobel", "noether",
	"northcutt", "noyce", "panini", "pare", "pasteur", "payne", "perlman", "pike",
	"poincare", "poitras", "ptolemy", "raman", "ramanujan", "ride", "ritchie", "rhodes",
	"roentgen", "rosalind", "rubin", "saha", "sammet", "shaw", "shirley", "shockley",
	"shtern", "sinoussi", "snyder", "spence", "stallman", "stonebraker", "swanson",
	"swartz", "swirles", "tesla", "tharp", "thompson", "torvalds", "tu", "turing",
	"varahamihira", "visvesvaraya", "volhard", "wescoff", "wilbur", "wiles", "williams",
	"williamson", "wilson", "wing", "wozniak", "wright", "wu", "yonath",
}

func GenerateRandomName() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	adj := adjectives[rng.Intn(len(adjectives))]
	name := surnames[rng.Intn(len(surnames))]

	return strings.ToLower(adj + "_" + name)
}
