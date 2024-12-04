package src

var CHANGE_TYPE = map[string]string{
	"FEAT":     "Adds or remove new feature",
	"FIX":      "Fixes a bug",
	"REFACTOR": "Restructure the code",
	"PERF":     "Improve Performance",
	"STYLE":    "White-space, formatting, etc.",
	"TEST":     "Add missing tests or correcting",
	"DOCS":     "Affect documentation only",
	"BUILD":    "Affect build components",
	"OPS":      "Affect operational components",
	"CHORE":    "Miscellaneous commits",
}

var GITMOJI_ARRAY = map[string]string{
	"✨":  ":sparkles: Introduce new features",
	"🎨":  ":art: Improve structure of the code",
	"⚡️": ":zap: Improve performance",
	"🔥":  ":fire: Remove code or files",
	"🐛":  ":bug: Fix a bug",
	"🚑️": ":ambulance: Critical hotfix",
	"📝":  ":memo: Add or update documentation",
	"🚀":  ":rocket: Deply stuff",
	"🎉":  ":tada: Begin a project",
	"✅":  ":white_check_mark: Add, update, or pass tests",
	"🔒️": ":lock: Fix security or privacy issue",
	"🔐":  ":closed_lock_with_key: Add or update secrets",
	"🔖":  ":bookmark: Release/Versions tags",
	"🚨":  ":rotating_light: Fix compiler/linter warnings",
	"🚧":  ":construction: Work in progress",
	"💚":  ":green_heart: Fix CI Build",
	"⬇️": ":arrow_down: Downgrade dependencies",
	"⬆️": ":arrow_up: Upgrade dependencies",
	"📌":  ":pushpin: Pin dependencies to specific version",
	"👷":  ":construction_worker: Add or update CI system",
	"♻️": ":recycle: Refactor code",
	"➕":  ":heavy_plus_sign: Add a dependency",
	"➖":  ":heavy_minus_sign: Remove a dependency",
	"🔧":  ":wrench: Add or update configuration files",
	"🔨":  ":hammer: Add or update development scripts",
	"✏️": ":pencil2: Fix typos",
	"💩":  ":poop: Write bad code that needs to be improved",
	"⏪️": ":rewind: Revert changes",
	"🔀":  ":twisted_rightwards_arrows: Merge branches",
	"📦️": ":package: Add or update compiled packages",
	"👽️": ":alien: Update code due to external API changes",
	"🚚":  ":truck: Move or renamed resources (paths, ...)",
	"💥":  ":boom: Introduce breaking changes",
	"🍱":  ":bento: Add or update assets",
	"🔊":  ":loud_sound: Add or update logs",
	"🔇":  ":mute: Remove logs",
	"🏗️": ":building_construction: Make architectural changes",
	"🤡":  ":clown_face: Mock things",
	"🙈":  ":see_no_evil: Add or update a .gitignore file",
	"⚰️": ":coffin: Remove dead code",
	"🧪":  ":test_tube: Add a failing test",
	"🧱":  ":bricks: Infrastructure releated changes",
	"🦺":  ":safety_vest: Add or update code for validation",
}

// String constants
const TITLE string = "CONVENTIONAL COMMITS CLI"
const TYPE string = "1. Select the type of change"
const GITMOJI string = "2. Select a gitmoji"
const MAIN_DESC string = "3. Write a Short Description"
const LONG_DESC string = "4. Write a Longer Description"
const VERSION string = "v0.1.0 - Riccardo La Marca"
const REPO string = "📦"
const BRANCH string = "🌲"
const REMOTE string = "👾"

const TITLE_Y int = 2
