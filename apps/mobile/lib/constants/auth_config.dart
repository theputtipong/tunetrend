const kGoogleServerClientId = String.fromEnvironment('GOOGLE_SERVER_CLIENT_ID');
const kGoogleIosClientId = String.fromEnvironment('GOOGLE_IOS_CLIENT_ID');

const kLoginEnabled = kGoogleServerClientId != '';

const kYouTubeDailySearchLimit = 10;
const kPlaylistNameMaxLength = 60;
