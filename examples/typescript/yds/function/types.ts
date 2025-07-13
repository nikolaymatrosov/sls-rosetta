export interface ProducerRequest {
  message: string;
  user_id: string;
  action: string;
}

export interface ProducerResponse {
  status_code: number;
  message: string;
  stream_id?: string;
}

export interface EventMetadata {
  event_id: string;
  event_type: string;
  created_at: string;
  cloud_id: string;
  folder_id: string;
}

export interface YDSDetails {
  stream_id: string;
  data: string;
}

export interface YDSMessage {
  event_metadata: EventMetadata;
  details: YDSDetails;
}

export interface YDSEvent {
messages: YDSMessage[];
}

export interface YDSResponse {
  status_code: number;
  message: string;
} 