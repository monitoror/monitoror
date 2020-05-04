export enum ConfigErrorId {
  CannotBeFetched = 'ERROR_CONFIG_CANNOT_BE_FETCHED',
  ConfigNotFound = 'ERROR_CONFIG_NOT_FOUND',
  ConfigVersionTooOld = 'ERROR_FIELD_TYPE_MISMATCH',
  FieldTypeMismatch = 'ERROR_FIELD_TYPE_MISMATCH',
  InvalidEscapedCharacter = 'ERROR_INVALID_ESCAPED_CHARACTER',
  InvalidFieldValue = 'ERROR_INVALID_FIELD_VALUE',
  MissingPathOrUrl = 'ERROR_MISSING_PATH_OR_URL',
  MissingRequiredField = 'ERROR_MISSING_REQUIRED_FIELD',
  UnauthorizedField = 'ERROR_UNAUTHORIZED_FIELD',
  UnauthorizedSubtileType = 'ERROR_UNAUTHORIZED_SUBTILE_TYPE',
  UnableToHydrate = 'ERROR_UNABLE_TO_HYDRATE',
  UnableToParseConfig = 'ERROR_UNABLE_TO_PARSE_CONFIG',
  UnexpectedError = 'ERROR_UNEXPECTED',
  UnknownField = 'ERROR_UNKNOWN_FIELD',
  UnknownNamedConfig = 'ERROR_UNKNOWN_NAMED_CONFIG',
  UnknownTileType = 'ERROR_UNKNOWN_TILE_TYPE',
  UnknownVariant = 'ERROR_UNKNOWN_VARIANT',
  UnsupportedVersion = 'ERROR_UNSUPPORTED_VERSION',
}

export default ConfigErrorId
