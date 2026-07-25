package repository

type evolutionOverride struct {
	Number  string
	Name    string
	Types   []string
	Trigger string
}

var pokemonAbilityOverrides = map[string][]string{
	"1":   {"Overgrow"},
	"2":   {"Overgrow"},
	"3":   {"Overgrow"},
	"4":   {"Blaze"},
	"5":   {"Blaze"},
	"6":   {"Blaze"},
	"7":   {"Torrent"},
	"8":   {"Torrent"},
	"9":   {"Torrent"},
	"15":  {"Swarm"},
	"25":  {"Static"},
	"35":  {"Cute Charm", "Magic Guard"},
	"51":  {"Sand Veil", "Arena Trap"},
	"95":  {"Rock Head", "Sturdy"},
	"108": {"Oblivious", "Own Tempo"},
	"109": {"Levitate"},
	"151": {"Synchronize"},
	"245": {"Pressure"},
	"306": {"Rock Head", "Sturdy"},
	"384": {"Air Lock"},
	"448": {"Inner Focus", "Steadfast"},
	"497": {"Overgrow"},
	"571": {"Illusion"},
	"609": {"Flash Fire", "Flame Body"},
	"613": {"Snow Cloak", "Slush Rush"},
	"733": {"Keen Eye", "Skill Link"},
}

var pokemonCategoryOverrides = map[string]string{
	"1":   "Seed",
	"2":   "Seed",
	"3":   "Seed",
	"4":   "Lizard",
	"5":   "Flame",
	"6":   "Flame",
	"7":   "Tiny Turtle",
	"8":   "Turtle",
	"9":   "Shellfish",
	"15":  "Poison Bee",
	"25":  "Mouse",
	"35":  "Fairy",
	"51":  "Mole",
	"95":  "Rock Snake",
	"108": "Licking",
	"109": "Poison Gas",
	"151": "New Species",
	"245": "Aurora",
	"306": "Iron Armor",
	"384": "Sky High",
	"448": "Aura",
	"497": "Regal",
	"571": "Illusion Fox",
	"609": "Luring",
	"613": "Chill",
	"733": "Cannon",
}

var pokemonTypeOverrides = map[string][]string{
	"15":  {"Inseto", "Venenoso"},
	"95":  {"Pedra", "Terrestre"},
	"306": {"Metal", "Terrestre"},
	"448": {"Lutador", "Metal"},
	"571": {"Noturno"},
	"609": {"Fantasma", "Fogo"},
	"733": {"Voador", "Normal"},
}

var pokemonEvolutionOverrides = map[string][]evolutionOverride{
	"1": {
		{Number: "1", Name: "Bulbasaur", Types: []string{"Grama", "Venenoso"}},
		{Number: "2", Name: "Ivysaur", Types: []string{"Grama", "Venenoso"}, Trigger: "Nível 16"},
		{Number: "3", Name: "Venusaur", Types: []string{"Grama", "Venenoso"}, Trigger: "Nível 36"},
	},
	"4": {
		{Number: "4", Name: "Charmander", Types: []string{"Fogo"}},
		{Number: "5", Name: "Charmeleon", Types: []string{"Fogo"}, Trigger: "Nível 16"},
		{Number: "6", Name: "Charizard", Types: []string{"Fogo", "Voador"}, Trigger: "Nível 36"},
	},
	"7": {
		{Number: "7", Name: "Squirtle", Types: []string{"Água"}},
		{Number: "8", Name: "Wartortle", Types: []string{"Água"}, Trigger: "Nível 16"},
		{Number: "9", Name: "Blastoise", Types: []string{"Água"}, Trigger: "Nível 36"},
	},
	"15": {
		{Number: "13", Name: "Weedle", Types: []string{"Inseto", "Venenoso"}},
		{Number: "14", Name: "Kakuna", Types: []string{"Inseto", "Venenoso"}, Trigger: "Nível 7"},
		{Number: "15", Name: "Beedrill", Types: []string{"Inseto", "Venenoso"}, Trigger: "Nível 10"},
	},
	"25": {
		{Number: "172", Name: "Pichu", Types: []string{"Elétrico"}},
		{Number: "25", Name: "Pikachu", Types: []string{"Elétrico"}, Trigger: "Nível de Amizade"},
		{Number: "26", Name: "Raichu", Types: []string{"Elétrico"}, Trigger: "Pedra do Trovão"},
	},
	"35": {
		{Number: "173", Name: "Cleffa", Types: []string{"Fada"}},
		{Number: "35", Name: "Clefairy", Types: []string{"Fada"}, Trigger: "Nível de Amizade"},
		{Number: "36", Name: "Clefable", Types: []string{"Fada"}, Trigger: "Pedra da Lua"},
	},
	"51": {
		{Number: "50", Name: "Diglett", Types: []string{"Terrestre"}},
		{Number: "51", Name: "Dugtrio", Types: []string{"Terrestre"}, Trigger: "Nível 26"},
	},
	"95": {
		{Number: "95", Name: "Onix", Types: []string{"Pedra", "Terrestre"}},
		{Number: "208", Name: "Steelix", Types: []string{"Metal", "Terrestre"}, Trigger: "Trocas"},
	},
	"108": {
		{Number: "108", Name: "Lickitung", Types: []string{"Normal"}},
		{Number: "463", Name: "Lickilicky", Types: []string{"Normal"}, Trigger: "Subir de Nível c/ Rollout"},
	},
	"109": {
		{Number: "109", Name: "Koffing", Types: []string{"Venenoso"}},
		{Number: "110", Name: "Weezing", Types: []string{"Venenoso"}, Trigger: "Nível 35"},
	},
	"151": {
		{Number: "151", Name: "Mew", Types: []string{"Psíquico"}},
	},
	"245": {
		{Number: "245", Name: "Suicune", Types: []string{"Água"}},
	},
	"306": {
		{Number: "304", Name: "Aron", Types: []string{"Metal", "Terrestre"}},
		{Number: "305", Name: "Lairon", Types: []string{"Metal", "Terrestre"}, Trigger: "Nível 32"},
		{Number: "306", Name: "Aggron", Types: []string{"Metal", "Terrestre"}, Trigger: "Nível 42"},
	},
	"384": {
		{Number: "384", Name: "Rayquaza", Types: []string{"Dragão"}},
	},
	"448": {
		{Number: "447", Name: "Riolu", Types: []string{"Lutador"}},
		{Number: "448", Name: "Lucario", Types: []string{"Lutador", "Metal"}, Trigger: "Nível de Amizade"},
	},
	"497": {
		{Number: "495", Name: "Snivy", Types: []string{"Grama"}},
		{Number: "496", Name: "Servine", Types: []string{"Grama"}, Trigger: "Nível 17"},
		{Number: "497", Name: "Serperior", Types: []string{"Grama"}, Trigger: "Nível 36"},
	},
	"571": {
		{Number: "570", Name: "Zorua", Types: []string{"Noturno"}},
		{Number: "571", Name: "Zoroark", Types: []string{"Noturno"}, Trigger: "Nível 30"},
	},
	"609": {
		{Number: "607", Name: "Litwick", Types: []string{"Fantasma", "Fogo"}},
		{Number: "608", Name: "Lampent", Types: []string{"Fantasma", "Fogo"}, Trigger: "Nível 41"},
		{Number: "609", Name: "Chandelure", Types: []string{"Fantasma", "Fogo"}, Trigger: "Pedra do Anoitecer"},
	},
	"613": {
		{Number: "613", Name: "Cubchoo", Types: []string{"Gelo"}},
		{Number: "614", Name: "Beartic", Types: []string{"Gelo"}, Trigger: "Nível 37"},
	},
	"733": {
		{Number: "731", Name: "Pikipek", Types: []string{"Voador", "Normal"}},
		{Number: "732", Name: "Trumbeak", Types: []string{"Voador", "Normal"}, Trigger: "Nível 14"},
		{Number: "733", Name: "Toucannon", Types: []string{"Voador", "Normal"}, Trigger: "Nível 28"},
	},
}

var pokemonDescriptionOverrides = map[string]string{
	"1":   "Há uma semente de planta nas costas desde o dia em que este Pokémon nasce. A semente cresce lentamente.",
	"2":   "Quando o bulbo nas costas cresce, parece perder a capacidade de ficar em pé nas patas traseiras.",
	"3":   "Sua planta floresce quando está absorvendo energia solar. Ele permanece em movimento para buscar a luz do sol.",
	"4":   "Tem preferência por coisas quentes. Quando chove, diz-se que o vapor jorra da ponta de sua cauda.",
	"5":   "Tem uma natureza bárbara. Na batalha, ele chicoteia sua cauda de fogo e corta com garras afiadas.",
	"6":   "Ele cospe fogo que é quente o suficiente para derreter pedregulhos. Pode causar incêndios florestais soprando chamas.",
	"7":   "Quando retrai seu longo pescoço em sua concha, esguicha água com força vigorosa.",
	"8":   "É reconhecido como um símbolo de longevidade. Se a concha tiver algas, esse Wartortle é muito antigo.",
	"9":   "Ele esmaga seu inimigo sob seu corpo pesado para causar desmaios. Em uma pitada, ele se retirará dentro de sua concha.",
	"15":  "Tem três ferrões venenosos nas patas dianteiras e na cauda. Eles são usados para espetar seu inimigo repetidamente.",
	"25":  "Pikachu, que pode gerar uma eletricidade poderosa, tem bolsas nas bochechas que são extra macias e super elásticas.",
	"35":  "Diz-se que a felicidade virá para aqueles que virem uma reunião de Clefairy dançando sob a lua cheia.",
	"51":  "Uma equipe de trigêmeos Diglett. Ele desencadeia enormes terremotos cavando 60 milhas no subsolo.",
	"95":  "À medida que escava o solo, absorve muitos objetos duros. Isso é o que torna seu corpo tão sólido.",
	"108": "Se a saliva pegajosa deste Pokémon entrar em contato com você e você não a limpar, uma coceira intensa se instalará. A coceira também não desaparecerá.",
	"109": "Seu corpo está cheio de gás venenoso. Ele flutua em lixões, procurando a fumaça do lixo cru e apodrecido.",
	"151": "Quando visto através de um microscópio, o cabelo curto, fino e delicado deste Pokémon pode ser visto.",
	"245": "Suicune encarna a compaixão de uma fonte de água pura. Ele atravessa a terra com graciosidade.",
	"306": "Aggron tem um chifre afiado o suficiente para perfurar grossas chapas de ferro. Ele derruba seus oponentes batendo neles primeiro com o chifre.",
	"384": "Diz-se que Rayquaza viveu por centenas de milhões de anos. Permanecem as lendas de como acabou o confronto entre Kyogre e Groudon.",
	"448": "Ele controla ondas conhecidas como auras, que são poderosas o suficiente para pulverizar rochas enormes. Ele usa essas ondas para derrubar sua presa.",
	"497": "Ele só dá tudo de si contra oponentes fortes que não se incomodam com o brilho dos olhos nobres de Serperior.",
	"571": "Este Pokémon se preocupa profundamente com outros de sua espécie e conjurará ilusões aterrorizantes para manter sua toca e sua mochila seguras.",
	"609": "Este Pokémon assombra mansões em ruínas. Ele balança seus braços para hipnotizar os oponentes com a dança sinistra de suas chamas.",
	"613": "Quando este Pokémon está bem de saúde, seu ranho fica mais grosso e pegajoso. Ele vai espalhar seu ranho em quem não gosta.",
	"733": "Eles batem bicos com outros de sua espécie para se comunicar. A força e o número de acertos dizem uns aos outros como eles se sentem.",
}
